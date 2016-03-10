//      DCSBT Index @InsZVA 2015
//      Based on WANG Sheng, QIN Xiaolin, SHEN Yao, et al. Research on durable CSB+-tree indexing technology. Journal of Frontiers of Computer Science and Technology, 2015, 9(2): 182-192.

package csbt

type CSBTNode struct {
    keyNum int
    lp, rp *CSBTNode
    key []uint
    isLeaf bool
}