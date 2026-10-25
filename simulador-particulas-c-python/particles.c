// particles.c
#include "particles.h"
#include <math.h>

// Atualiza a posição de uma partícula
void update_particle(Particle* p, double dt) {
    p->x += p->vx * dt;
    p->y += p->vy * dt;
}

// Aplica a força gravitacional entre duas partículas
void apply_gravity(Particle* p1, Particle* p2, double dt) {
    double G = 6.67430e-11;  // Constante gravitacional
    double dx = p2->x - p1->x;
    double dy = p2->y - p1->y;
    double distance = sqrt(dx * dx + dy * dy);
    double force = G * p1->mass * p2->mass / (distance * distance);
    double fx = force * dx / distance;
    double fy = force * dy / distance;

    p1->vx += fx / p1->mass * dt;
    p1->vy += fy / p1->mass * dt;
    p2->vx -= fx / p2->mass * dt;
    p2->vy -= fy / p2->mass * dt;

    p1->x += p1->vx * dt;
    p1->y += p1->vy * dt;
    p2->x += p2->vx * dt;
    p2->y += p2->vy * dt;
}

// Verifica e resolve colisão entre duas partículas
int check_collision(Particle* p1, Particle* p2) {
    double dx = p2->x - p1->x;
    double dy = p2->y - p1->y;
    double dist = sqrt(dx * dx + dy * dy);
    double min_dist = 1.0; // Distância mínima para colisão (raio somado das partículas)

    if (dist < min_dist) {
        double overlap = min_dist - dist;
        p1->x -= dx / dist * overlap / 2;
        p1->y -= dy / dist * overlap / 2;
        p2->x += dx / dist * overlap / 2;
        p2->y += dy / dist * overlap / 2;

        // Inverte as velocidades
        double temp_vx = p1->vx;
        double temp_vy = p1->vy;
        p1->vx = p2->vx;
        p1->vy = p2->vy;
        p2->vx = temp_vx;
        p2->vy = temp_vy;

        return 1;
    }
    return 0;
}
