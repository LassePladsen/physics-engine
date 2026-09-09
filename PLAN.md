This is a plan made by chatGPT just in case I feel lost in terms of what to focus on next.

# Physics Engine Plan

Build a simple 2D physics engine incrementally. Keep each step small, testable, and visual.

## 1. Math Foundation

* [x] Implement `Vec2`
* [x] Add: `Add`, `Sub`, `Mul`, `Dot`, `Length`, `Normalize`
* [x] Add unit tests

## 2. Particles

* [x] Implement `Particle`
* [x] Position, velocity, mass
* [ ] Implement basic time integration: Update()
* [x] Verify constant-velocity motion
* [ ] Implement acceleration. Should this be stored in the Particle struct or be retrieved in Update()?

## 3. World

* [ ] Implement `World`
* [ ] Store particles
* [ ] Implement `World.Step(dt)`
* [ ] Add configurable gravity
* [ ] Verify free-fall motion

## 4. Visualization

* [ ] Create simple 2D visualization
* [ ] Render particles
* [ ] Run simulation loop
* [ ] Visualize gravity and motion

## 5. Boundaries

* [ ] Add world/floor boundaries
* [ ] Detect particle-boundary collision
* [ ] Resolve penetration
* [ ] Add restitution
* [ ] Visualize bouncing particles

## 6. Particle Collisions

* [ ] Add particle radius
* [ ] Detect particle-particle collisions
* [ ] Resolve collisions
* [ ] Preserve momentum
* [ ] Add restitution

## 7. Forces

* [ ] Add generic forces
* [ ] Implement drag
* [ ] Implement friction
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

