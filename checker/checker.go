package checker

import (
	"barker-go/vcr/message"
	"net/http"
)
type checker interface {
	updateConfig()
	resolve(request *message.BidReq) *FinalResult
}

const (
	KeyOfBucket = "model_various/augmentor_config.json"
	Resolved    = true
	Unresolved  = false
)

type FinalResult struct {
	IsResolved bool
	IsBidding bool
	SegmentList []string // will be changed to []Segment
}

type AugmentorResponseChecker struct {
	region string
	bucket string
	next *AugmentorResponseChecker
	checker
}

func (checker *AugmentorResponseChecker) updateFilePath(region string, bucket string) {
	checker.region = region
	checker.bucket = bucket
	checker.next.updateFilePath(region,bucket)
}

func (checker *AugmentorResponseChecker) setNext(next *AugmentorResponseChecker) {
	checker.next= next
}

func (checker *AugmentorResponseChecker) getSegment(bidRequest *message.BidReq) *FinalResult{
	result := checker.resolve(bidRequest)

	if result.IsResolved {
		return result
	}else if checker.next!=nil {
		return checker.next.getSegment(bidRequest)
	}else{
		return &FinalResult{false,false,make([]string,0)}
	}
}
