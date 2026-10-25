struct TreeNode* invertTree(struct TreeNode* root){
    if (root != NULL){
        invertTree(root->right);
        invertTree(root->left);
        struct TreeNode* temp = root->right;
        root->right = root->left;
        root->left = temp;
    }

    return root;

}