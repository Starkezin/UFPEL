int findSecondMinimumValue(struct TreeNode* root){
    struct TreeNode* sentinela;
    sentinela = root;
    int valraiz, menEsq, menDir;
    if(root->left == NULL){
        return -1;
    }

    valraiz = root->val;
    root = root->left;
    menEsq = root->val;
    root = sentinela;
    root = root->right;
    menDir = root->val;
    if(valraiz == menEsq && valraiz == menDir){
        return -1;
    }
    else if(valraiz <= menEsq && valraiz >= menDir){
        return valraiz;
    }
    else if(menEsq <= valraiz && menEsq >= menDir){
        return menEsq;
    }
    else if(menEsq == menDir && valraiz > menEsq){
        return valraiz;
    }
    else if(menDir == valraiz && menEsq > menDir){
        return menEsq;
    }
    else{
        return menDir;
    }

    return root->val;
}