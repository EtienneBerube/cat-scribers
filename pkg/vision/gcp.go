package vision

import (
	"cloud.google.com/go/vision/apiv1"
	"context"
	"encoding/base64"
	"fmt"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"log"
)

func HasCat(imgB64 string, name string) (bool, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	bytes, _ := base64.StdEncoding.DecodeString(imgB64)

	res, err := client.AnnotateImage(ctx, &pb.AnnotateImageRequest{
		Image: &pb.Image{
			Content: bytes,
		},
		Features: []*pb.Feature{
			{
				Type:       pb.Feature_OBJECT_LOCALIZATION,
				MaxResults: 10,
			},
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	if res.GetError() != nil {
		log.Fatal(res.GetError().Message)
	}

	if len(res.LocalizedObjectAnnotations) == 0 {
		fmt.Println("No Objects found.")
	} else {
		fmt.Printf("Image: %s, Objects:\n", name)
		for _, object := range res.LocalizedObjectAnnotations {
			fmt.Printf("%s\n", object.Name)
			if object.Name == "Cat" {
				return true, nil
			}
		}
	}

	return false, nil
}

func HasCatMultiple(imgsB64 []string, names []string) ([]bool, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	var reqs []*pb.AnnotateImageRequest

	for _, img := range imgsB64 {
		bytes, _ := base64.StdEncoding.DecodeString(img)
		reqs = append(reqs, &pb.AnnotateImageRequest{
			Image: &pb.Image{
				Content: bytes,
			},
			Features: []*pb.Feature{
				{
					Type:       pb.Feature_OBJECT_LOCALIZATION,
					MaxResults: 10,
				},
			},
		})
	}

	responses, err := client.BatchAnnotateImages(ctx, &pb.BatchAnnotateImagesRequest{
		Requests: reqs,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	imageHasCat := make([]bool, len(responses.Responses))

	for i, resp := range responses.Responses {
		if resp.GetError() != nil {
			imageHasCat[i] = false
		}
		if len(resp.LocalizedObjectAnnotations) == 0 {
			fmt.Printf("Image: %s, No Objects found\n", names[i])
		} else {
			fmt.Printf("Image: %s, Objects:\n", names[i])
			for _, object := range resp.LocalizedObjectAnnotations {
				fmt.Printf("%s\n", object.Name)
				if object.Name == "Cat" {
					imageHasCat[i] = true
					break
				}
			}
			imageHasCat[i] = false
		}
	}

	return imageHasCat, nil
}
