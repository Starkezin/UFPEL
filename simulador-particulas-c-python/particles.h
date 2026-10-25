// particles.h
#ifndef PARTICLES_H
#define PARTICLES_H

typedef struct {
    double x, y;
    double vx, vy;
    double mass;
} Particle;

void update_particle(Particle* p, double dt);
void apply_gravity(Particle* p1, Particle* p2, double dt);
int check_collision(Particle* p1, Particle* p2);

#endif
