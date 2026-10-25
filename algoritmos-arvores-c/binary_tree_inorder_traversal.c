int count(struct TreeNode* root)
{
    if(root)
        return 1 + count(root->left) + count(root->right);
    return 0;
}

void sort(struct TreeNode* raiz, int** res)
{
    if(raiz)
    {
        sort(raiz->left, res);
        *((*res)++) = raiz->val;
        sort(raiz->right, res);
    }
}


int* inorderTraversal(struct TreeNode* root, int* returnSize)
{
    *returnSize = count(root);
    int* res = (int*)malloc(((*returnSize)) * sizeof(int));
    memset(res, 0, (*returnSize) * sizeof(int));   
    int *cpy = res;
    sort(root, &cpy);
    return res;
}