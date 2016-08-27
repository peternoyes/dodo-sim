# dodo-sim
6502 Simulator for Dodo the 6502 Portable Game System

[![Build Status](https://travis-ci.org/peternoyes/dodo-sim.svg?branch=master)](https://travis-ci.org/peternoyes/dodo-sim)

dodo-sim is the core simulator for [Dodo](https://github.com/peternoyes/dodo) which is a 6502 based homebrew game system. The simulator is for the most part a port of http://rubbermallet.org/fake6502.c with some 65C02 opcodes. Decimal mode is also fixed. The simulator passes the Klaus set of 6502 tests found [here](https://github.com/Klaus2m5)

This is a library that must be used by an outside application to initiate the simulator. Call dodosim.Simulate to run. That function takes in a user implemented Renderer.

The simulator is hardcoded to use the address space layout and devices used in Dodo but it could be repurposed for other systems. 

## Todo

- Sound simulation