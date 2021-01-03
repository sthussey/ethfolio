# ethfolio

ethfolio attempts to extract information from the Ethereum
blockchain in order to build a FIFO portfolio to calculate
earnings and capital gains. The goal is recognize several
types of income types without being opinionated on the
tax implications.

## Income

Initial MVP is accumulating income transactions recognized
by the source or destination wallet address. Such a system
could be used to track payments from a mining pool for example.

## Staking

The Ethereum 2.0 proof-of-stake network has been launched and
it is not clear if rewards for those staking validators on the
network must claim be claimed even before they can be extracted.

## Capital Gains

ethfolio will function as a FIFO lot tracker for Eth transactions. With
a highly volatile exchange rate between Eth and sovereign currencies,
calculating capital gains or lesses is non-trivial with more than a
handful of transactions. Below is the method for tracking your ethfolio.

  1. You must provide exchange rate history. Currently the expectation is to
     track income at a daily granularity. All income transactions will be accumulated
     for a given day and attributed the daily exchange rate.
  2. The ethfolio will be stored locally to avoid re-calculating blockchain
     transactions. This also assumes additional income sources are
     not specified prior to the latest date in the local store.
  3. When you make exchanges, you specify the actual exchange rate and date
     and that exchange will be tracked as a unique lot without being
     co-mingled with income lots.
    
