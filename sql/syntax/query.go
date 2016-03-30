package syntax

func Parser(tr *TokenReader) (*syntaxTreeNode, error) {
    fork := tr.Fork()
    node, err := createtableParser(fork)
    // node, err := filedParser(tr)
    // if err == nil {
    //     return node, nil
    // }
    return node, err
}