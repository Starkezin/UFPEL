void busca(struct Node* root, int* res, int* count) {
    if (!root) return;
    for (int i = 0; i < root->numChildren; i++) {
        busca(root->children[i], res, count);
    }
    res[(*count)++] = root->val;
}

int* postorder(struct Node* root, int* returnSize) {
    int* res = (int*) malloc(10000 * sizeof(int));
    int count = 0;
    busca(root, res, &count);
    *returnSize = count;
    return res;

}