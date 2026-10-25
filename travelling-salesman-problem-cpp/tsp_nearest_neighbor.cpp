#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <limits>
#include <chrono>
using namespace std;

vector<int> convertVecStrToVecInt(const vector<string>& vecStrings) 
{
    vector<int> vecInt;
    for (auto& str : vecStrings) 
    {
        int aux = atoi(str.c_str());
        vecInt.push_back(aux);
    }
    return vecInt;
}
vector<string> getTokens(const string& line) 
{
    vector<string> tokens;
    string token;
    for (char c : line) 
    {
        if (c != ' ') token += c;
        else if (!token.empty()) 
        {
            tokens.push_back(token);
            token.clear();
        }
    }
    if (!token.empty()) tokens.push_back(token);
    return tokens;
}
int calcTSP(const vector<int>& route, const vector<vector<int>>& distances) 
{
    int total_distance = 0;
    for (size_t i = 0; i < route.size() - 1; ++i)
        total_distance += distances[route[i]][route[i + 1]];
    total_distance += distances[route.back()][route[0]];
    return total_distance;
}
vector<int> nearestNeighbor(const vector<vector<int>>& distances, int start = 0) 
{
    int n = distances.size();
    vector<bool> visited(n, false);
    vector<int> route;
    route.reserve(n);
    route.push_back(start);
    visited[start] = true;

    int current = start;
    for (int step = 1; step < n; ++step) 
    {
        int next_city = -1;
        int min_dist = numeric_limits<int>::max();

        for (int j = 0; j < n; ++j) 
        {
            if (!visited[j] && distances[current][j] < min_dist) 
            {
                min_dist = distances[current][j];
                next_city = j;
            }
        }
        route.push_back(next_city);
        visited[next_city] = true;
        current = next_city;
    }

    return route;
}

int main(int argc, char **argv) 
{
    string filename = argv[1];
    ifstream fileInput(filename);
    string line;
    vector<vector<int>> distances;

    if (!fileInput) 
    {
        cerr << "Erro ao abrir arquivo!" << endl;
        return 1;
    }

    while (getline(fileInput, line)) 
    {
        if (line.empty()) continue;
        vector<string> auxString = getTokens(line);
        vector<int> auxInt = convertVecStrToVecInt(auxString);
        distances.push_back(auxInt);
    }

    auto inicio = std::chrono::high_resolution_clock::now();
    vector<int> route = nearestNeighbor(distances, 0);
    int distance = calcTSP(route, distances);

    auto fim = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::milli> duracao = fim - inicio;
    cout << "Tempo de execucao: " << duracao.count() << " ms" << endl;
    
    cout << "Rota aproximada: ";
    for (int city : route) 
    {
        cout << city << " ";
    }   
    cout << endl;


    cout << "Distancia total aproximada: " << distance << endl;
    return 0;
}
