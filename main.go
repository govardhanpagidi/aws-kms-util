package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func main() {
	regions := []string{
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"sa-east-1",
		"me-south-1",
		"af-south-1",
		"eu-south-1",
		"ap-southeast-3",
		"ap-east-1",
	}

	for i := range regions {

		sess, err := session.NewSession(&aws.Config{Region: aws.String(regions[i])})
		if err != nil {
			fmt.Println(err)
			return
		}

		svc := kms.New(sess)

		max := int64(100)
		listInput := &kms.ListKeysInput{
			Limit: &max,
		}

		result, err := svc.ListKeys(listInput)
		if err != nil {
			fmt.Printf("error inlisting secrets: %+v", err.Error())
			return
		}
		for i := range result.Keys {

			fmt.Println("")

			dscKey := &kms.DescribeKeyInput{KeyId: result.Keys[i].KeyId}
			key, err := svc.DescribeKey(dscKey)
			if err != nil {
				fmt.Println(err)
				break
			}

			if *key.KeyMetadata.KeyState == "Disabled" {
				fmt.Printf("%+v", key)
				fmt.Printf("already disabled : %s", key.String())

			}
			if *key.KeyMetadata.KeyManager == "CUSTOMER" && *key.KeyMetadata.KeyState == "Enabled" {
				fmt.Printf("%+v", key)

				_, err := svc.DisableKey(&kms.DisableKeyInput{KeyId: key.KeyMetadata.KeyId})

				// Deleting keys. may not work check once.
				//_, err := svc.DeleteCustomKeyStore(&kms.DeleteCustomKeyStoreInput{CustomKeyStoreId: key.KeyMetadata.CustomKeyStoreId})

				if err != nil {
					fmt.Printf("error in deleting secret: %+v", err.Error())
				}
				fmt.Printf("Disable key : %s in region %s", *key.KeyMetadata.KeyId, regions[i])
			}
		}
	}

}
