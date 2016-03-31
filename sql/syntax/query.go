package syntax

func Parser(tr *TokenReader) (*SyntaxTreeNode, error) {
    fork := tr.Fork()
    node, err := createtableParser(fork)
    // node, err := filedParser(tr)
    // if err == nil {
    //     return node, nil
    // }
    return node, err
}