- [x] Math
- [x] particle
- [x] visulization
- [x] Bounce in bounds
- [x] Global gravity: maybe add acceleration logic in particle?
- [x] array of particles
- [ ] particle-to-particle gravity
- [ ] simulate a star system

Below is a plan made by chatGPT just in case I feel lost in terms of what to focus on next.

# Physics Engine Plan

Build a simple 2D physics engine incrementally. Keep each step small, testable, and visual.

## 1. Math Foundation

* [x] Implement `Vec2`
* [x] Add: `Add`, `Sub`, `Mul`, `Dot`, `Length`, `Normalize`
* [x] Add unit tests

## 2. Particles

* [x] Implement `Particle`
* [x] Position, velocity, mass
* [x] Implement basic time integration: Update()
* [x] Verify constant-velocity motion
* [x] Implement acceleration. Should this be stored in the Particle struct or be retrieved in Update()?

## 3. Visualization

* [x] Create simple 2D visualization
* [x] Render particles
* [x] Run simulation loop
* [x] Visualize gravity and motion

## 4. World

* [x] Implement `World`
* [x] Store particles
* [x] Implement `World.Step(dt)`
* [x] Add configurable gravity
* [x] Verify free-fall motion


## 5. Boundaries

* [x] Add world/floor boundaries
* [x] Detect particle-boundary collision
* [x] Resolve penetration
* [x] Add restitution
* [x] Visualize bouncing particles

## 6. Particle Collisions

* [x] Add particle radius
* [x] Detect particle-particle collisions
* [x] Resolve collisions
* [x] Preserve momentum
* [x] Add restitution

## 7. Forces

* [ ] Add generic forces
* [ ] Implement drag
* [x] Implement friction
* [ ] Implement impulses
* [ ] Add springs

## 8. Rigid Bodies

* [ ] Separate particles from rigid bodies
* [ ] Add orientation and angular velocity
* [ ] Add torque
* [ ] Add moment of inertia
* [ ] Implement rigid-body collision detection/resolution

## 9. Collision System

* [ ] Define shapes/colliders
* [ ] Circle collisions
* [ ] AABB collisions
* [ ] Polygon collisions
* [ ] Add contact manifolds
* [ ] Add broad-phase detection

## 10. Constraints

* [ ] Distance constraints
* [ ] Springs
* [ ] Joints
* [ ] Constraint solver

## Guiding Principles

* Keep the engine dependency-light.
* Prefer simple implementations before optimized ones.
* Add visualization whenever a feature is implemented.
* Add tests for math and physics behavior.
* Don't implement future features prematurely.
* Each milestone should leave the engine in a working state.

