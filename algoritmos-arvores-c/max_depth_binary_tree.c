int maxDepth(struct TreeNode* root){
    if (root == NULL){
        return 0;
    }
    int leftDeep = 1, rightDeep = 1;
    leftDeep += maxDepth(root->left);
    rightDeep += maxDepth(root->right);

    if(rightDeep > leftDeep) return rightDeep;
    else return leftDeep;
}
