package decisions

// For returns the (apply, review) pair the engine bands a purpose's answers
// with. Noul answers are banded on P(yes); choice and score answers on
// confidence. Callers that gate on another quantity (P(off_topic), the sum of
// kept roles, the primary-concept probability, ...) call BandFor with the
// named threshold they need.
//
//	relevance       relevance_quarantine / relevance_review
//	quality         quality_apply / quality_review
//	duplicate       duplicate / duplicate
//	classify        0.5 / 0.5 (facets are written from the answer itself)
//	concept         concept_noul / primary_concept
//	link            link_apply / link_review
//	mention         mention_sense / link_review
//	find            find_keep / find_keep
//	okf_type        okf_type / 0.5
//	contract_check  contract_conflict / contract_conflict
//	anything else   0.8 / 0.5
func (t Thresholds) For(purpose Purpose) (apply, review float64) {
	switch purpose {
	case PurposeRelevance:
		return t.Get("relevance_quarantine"), t.Get("relevance_review")
	case PurposeQuality:
		return t.Get("quality_apply"), t.Get("quality_review")
	case PurposeDuplicate:
		return t.Get("duplicate"), t.Get("duplicate")
	case PurposeClassify:
		return 0.5, 0.5
	case PurposeConcept:
		return t.Get("concept_noul"), t.Get("primary_concept")
	case PurposeLink:
		return t.Get("link_apply"), t.Get("link_review")
	case PurposeMention:
		return t.Get("mention_sense"), t.Get("link_review")
	case PurposeFind:
		return t.Get("find_keep"), t.Get("find_keep")
	case PurposeOKFType:
		return t.Get("okf_type"), 0.5
	case PurposeContractCheck:
		return t.Get("contract_conflict"), t.Get("contract_conflict")
	default:
		return 0.8, 0.5
	}
}

// band sets Answer.Band for a decided answer from the purpose's pair.
func band(answer Answer, apply, review float64) Answer {
	if !answer.Decided() {
		answer.Band = ""
		return answer
	}
	switch {
	case answer.Noul != nil:
		answer.Band = BandFor(*answer.Noul, apply, review)
	case answer.Confidence != nil:
		answer.Band = BandFor(*answer.Confidence, apply, review)
	default:
		answer.Band = BandIgnore
	}
	return answer
}
