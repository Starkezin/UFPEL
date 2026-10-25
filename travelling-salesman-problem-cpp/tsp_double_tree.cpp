#include <iostream>
#include <vector>
#include <fstream>
#include <algorithm>
#include <sstream>
#include <string>
#include <chrono>
using namespace std;

const int INF = 10000000; // weight INF means there is no edge
int n;
vector<vector<int>> adj; // adjacency matrix of graph
vector<vector<int>> mst;
vector<bool> visited;

struct Edge {
    int w = INF, to = -1;
};


vector<int> convertVecStrToVecInt(vector<string> vecStrings)
{
    vector<int> vecInt;
    for (auto str : vecStrings)
    {
        int aux = atoi(str.c_str());
        vecInt.push_back(aux);
    }
    return vecInt;
}
int calcTSP(vector<int> route, vector<vector<int>> distances)
{
    int total_distance = 0;
    for (int i = 0; i < route.size() - 1; i++)
    {
        total_distance += distances[route[i]][route[i + 1]];
    }
    total_distance += distances[route[route.size() - 1]][route[0]];
    return total_distance;
}
vector<string> getTokens(string line)
{
    int pos = 0;
    string token;
    vector<string> vecStrings;
    for (pos = 0; pos < line.size(); ++pos)
    {
        if (line[pos] != ' ')
        {
            token.insert(token.end(), line[pos]);
        }
        else
        {
            if (!token.empty())
            {
                vecStrings.push_back(token);
                token = "";
            }
        }
    }
    vecStrings.push_back(token);
    return vecStrings;
}

vector<vector<int>> final_vec;
void prim() {
    int total_weight = 0;
    vector<bool> selected(n, false);
    vector<Edge> min_e(n);
    
    min_e[0].w = 0;

    mst.assign(n, {});

    for (int i=0; i<n; ++i) {
        int v = -1;
        for (int j = 0; j < n; ++j) {
            if (!selected[j] && (v == -1 || min_e[j].w < min_e[v].w))
                v = j;

        }

        if (min_e[v].w == INF) {
            cout << "No MST!" << endl;
            exit(0);
        }

        selected[v] = true;
        if (min_e[v].to != -1){
            int u = min_e[v].to;
            mst[u].push_back(v);
            mst[v].push_back(u);
        }

        for (int to = 0; to < n; ++to) {
            if (adj[v][to] < min_e[to].w)
                min_e[to] = {adj[v][to], v};
        }
    }

}

void dfs(int v, vector<int> &order){
    visited[v] = true;
    order.push_back(v);
    for(auto to:mst[v]){
        if(!visited[to]){
            dfs(to, order);
            order.push_back(v);
        }
    }
}

int calcRouteCost(const vector<int> &route) {
    int cost = 0;
    for (size_t i = 0; i < route.size() - 1; ++i) {
        cost += adj[route[i]][route[i + 1]];
    }
    cost += adj[route.back()][route[0]];
    return cost;
}

int main(int argc, char** argv) {
    string filename = argv[1];
    ifstream file(filename);
    if (!file) { cout << "Erro abrindo arquivo!\n"; return 1; }

    string line;
    while (getline(file, line)) {
        vector<int> row;
        int x;
        stringstream ss(line);
        while (ss >> x) row.push_back(x);
        adj.push_back(row);
    }
    auto ini = chrono::high_resolution_clock::now();

    n = adj.size();
    prim();

    visited.assign(n, false);
    vector<int> walk;
    dfs(0, walk); 

    vector<int> seen(n, false), hamiltonian;
    for (int v : walk)
        if (!seen[v]) { hamiltonian.push_back(v); seen[v] = true; }

    int cost = calcRouteCost(hamiltonian);
    auto fim = chrono::high_resolution_clock::now();
    chrono::duration<double, milli> duracao = fim - ini;
    cout << "Tempo de execucao (ms): " << duracao.count() << endl;

    cout << "Caminho aproximado: ";
    for (int v : hamiltonian) cout << v << " ";
    cout << "\nCusto total: " << cost << endl;
}