import ctypes
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation

# Carregue a biblioteca
import os
import ctypes

current_dir = os.path.dirname(os.path.abspath(__file__))
dll_path = os.path.join(current_dir, 'particles.so')
particles_lib = ctypes.CDLL(dll_path)
# Defina a classe Particle
class Particle(ctypes.Structure):
    _fields_ = [("x", ctypes.c_double),
                ("y", ctypes.c_double),
                ("vx", ctypes.c_double),
                ("vy", ctypes.c_double),
                ("mass", ctypes.c_double)]

# Defina a função apply_gravity
apply_gravity = particles_lib.apply_gravity
apply_gravity.argtypes = [ctypes.POINTER(Particle), ctypes.POINTER(Particle), ctypes.c_double]
apply_gravity.restype = None

# Crie partículas para o exemplo
num_particles = 15
particles = (Particle * num_particles)()

# Inicialize partículas (exemplo)
for i in range(num_particles):
    particles[i].x = np.random.rand()
    particles[i].y = np.random.rand()
    particles[i].vx = np.random.rand()
    particles[i].vy = np.random.rand()
    particles[i].mass = 1.0

# Função de atualização
def update(frame):
    dt = 0.1
    for i in range(num_particles):
        for j in range(i + 1, num_particles):
            apply_gravity(ctypes.byref(particles[i]), ctypes.byref(particles[j]), dt)
    # Atualize a visualização (exemplo básico)
    plt.clf()
    x = [p.x for p in particles]
    y = [p.y for p in particles]
    plt.scatter(x, y)

# Configuração da animação
fig, ax = plt.subplots()
ani = FuncAnimation(fig, update, frames=range(100), interval=100)

plt.show()