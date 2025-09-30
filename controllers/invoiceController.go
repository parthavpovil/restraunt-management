package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthav/restraunt-management/database"
	"github.com/parthav/restraunt-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceViewFormat struct{
	Invoice_id		string
	Payment_method		string
	Order_id		string
	Payment_status		*string
	Payment_due		interface{}
	Table_number	 interface{}
	Payment_due_date		time.Time
	Order_details		interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client,"invoices")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx , cancel :=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		result, err :=invoiceCollection.Find(ctx,bson.M{})
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while listing invoices",
			})
			return 
		}
		var allInvoices []bson.M

		err=result.All(ctx,&allInvoices)
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while listing invoices",
			})
			return 
		}
		c.JSON(http.StatusOK,allInvoices)

	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel =context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		invoiceID :=c.Param("invoice_id")
		var invoice models.Invoice
		err :=invoiceCollection.FindOne(ctx,bson.M{"invoice_id":invoiceID}).Decode(&invoice)
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while retriving invoice",
			})
			return 
		}
		var invoiceview InvoiceViewFormat
		allOrderItems, err:=ItemsByOrder(invoiceID)
		invoiceview.Order_id =invoice.Order_id
		invoiceview.Payment_due_date=invoice.Payment_due_date
		invoiceview.Payment_method="null"

		if invoice.Payment_method !=nil{
			invoiceview.Payment_method= *invoice.Payment_method
		}

		invoiceview.Invoice_id =invoice.Invoice_id
		invoiceview.Payment_status =*&invoice.Payment_status
		invoiceview.Payment_due=allOrderItems[0]["payment_due"]
		invoiceview.Table_number=allOrderItems[0]["table_number"]
		invoiceview.Order_details=allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK,invoiceview)
		


	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var invoice models.Invoice
		err :=c.BindJSON(&invoice)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),			})
			return 
		}
		validationErr :=validate.Struct(invoice)
		if validationErr !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":validationErr.Error(),})
				return 
		}
		var order models.Order
		err =orderCollection.FindOne(ctx,bson.M{"order_id":invoice.Order_id}).Decode(&order)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"order ws not found",
			})
			return 
		}
		status :="PENDING"

		if invoice.Payment_status ==nil{
			invoice.Payment_status =&status
		}
		invoice.Payment_due_date =time.Now().AddDate(0,0,1)
		invoice.Created_at=time.Now()
		invoice.Updated_at=time.Now()
		invoice.ID=primitive.NewObjectID()
		invoice.Invoice_id=invoice.ID.Hex()

		result, err :=invoiceCollection.InsertOne(ctx,invoice)
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error while adding invoice",
			})
			return 
		}
		c.JSON(http.StatusOK,result)

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceID:=c.Param("invoice_id")
		err :=c.BindJSON(&invoice)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),			})
			return 
		}

		filter := bson.M{"invoice_id":invoiceID}

		 var updateObj primitive.D

		if invoice.Payment_method !=nil{
			updateObj =append(updateObj, bson.E{"payment_method",invoice.Payment_method})
		}

		if invoice.Payment_status !=nil{
			updateObj=append(updateObj, bson.E{"payment_status",invoice.Payment_status})
		}

		invoice.Updated_at =time.Now()
		updateObj=append(updateObj, bson.E{"updated_at", invoice.Updated_at})

		upsert :=true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status :="PENDING"

		if invoice.Payment_status ==nil{
			invoice.Payment_status =&status
		}

		result, err :=invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"invoice update failed",
			})
			return 
		}
		c.JSON(http.StatusOK,result)
	}
}
