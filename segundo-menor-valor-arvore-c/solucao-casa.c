int findSecondMinimumValue(struct TreeNode* root) {
    if (root == NULL) return -1;

    int menor = -1;
    int segundomenor = -1;

    void busca(struct TreeNode* root) {
        if (root == NULL) return;

        if(root->val < menor || menor == -1){
            segundomenor = menor;
            menor = root->val;
        }

        else if ( (root->val < segundomenor || segundomenor == -1 )&& root->val != menor) segundomenor = root->val;

        busca(root->left);
        busca(root->right);
    }

    busca(root);
    
    return segundomenor;
}
