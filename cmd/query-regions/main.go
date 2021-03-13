package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"fmt"
	"github.com/ashwanthkumar/golang-utils/sets"
	"log"
	"math/rand"
	"os"
)

// RegQuery holds the regions.
// once it is filled the first time it only keeps the matching items
// [].Add([a,b,c,d]) -> [a,b,c,d]
// [a,b,c,d].Add([c,d]) -> [c,d]
type RegQuery struct {
	regions sets.Set
	sess    *session.Session
}

// NewRegQuery returns an empty rq
func NewRegQuery(sess *session.Session) RegQuery {
	return RegQuery{
		regions: sets.Empty(),
		sess:    sess,
	}
}

func (rq *RegQuery) queryRegions(inst string) (err error) {
	debug := false
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	regs := sets.Empty()
	svc := ec2.New(rq.sess)
	input := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: aws.String("availability-zone"),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []*string{aws.String(inst)},
			},
		},
	}
	offers, err := svc.DescribeInstanceTypeOfferings(input)
	if err != nil {
		panic(err)
	}
	for _, offer := range offers.InstanceTypeOfferings {
		regs.Add(*offer.Location)
	}
	if rq.regions.Size() == 0 {
		if debug {
			log.Printf("Empty regions; set: %v", regs.Values())
		}
		rq.regions = regs
	} else {
		if debug {
			log.Printf("Current regions: %v", rq.regions.Values())
			log.Printf("New regions: %v", regs.Values())
		}
		rq.regions = rq.regions.Intersect(regs)
		if debug {
			log.Printf("Intersect regions: %v", rq.regions.Values())
		}
	}
	return
}

func (rq RegQuery) String() (res []string) {
	for _, ele := range rq.regions.Values() {
		res = append(res, ele)
	}
	return
}

func (rq RegQuery) RandomPick() {
	in := rq.regions.Values()
	randomIndex := rand.Intn(len(in))
	pick := in[randomIndex]
	fmt.Println(pick)
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		panic(err)
	}
	rq := NewRegQuery(sess)
	// Create EC2 service client
	for _, inst := range os.Args[1:] {
		rq.queryRegions(inst)
	}
	rq.RandomPick()
}

