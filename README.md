# DistributedMutualExclusion

eepy consensus.

To run and showcase mutual exclusion:

1. Open three seperate terminals
2. In the first terminal enter 'go run . \<X>\' where X is the last digit of the 5000s port you want to use.
3. Do the same in a second terminal, instead entering 'go run . \<X>+1\', X+1 being the first port + 1.
4. Do the same in a third terminal, instead entering 'go run . \<X>+2\', X+1 being the first port + 2.
5. The given peers in each terminal should originally go back and forth.

Credit:
Peer-to-Peer architecture template from: https://github.com/NaddiNadja/peer-to-peer
