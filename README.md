# esports-stock-market


why?

because current esports fantasy leagues suck.
because they are all independent to each esport....when you could share code/systems across all

how make better?

main thing is to have the usual x-budget which you 'buy' players for your team with.
However to make more fun/skilly, the price of players is dynamic throughout the season, rising and falling based on popularity.
So it's like a stock market.

one other fantasy league idea I like is being automaticlly entered into a league with random players, ON TOP OF having friend leagues. Not everyone has enough friends interested to make a good league.

implementation?

using go and postgres for web shenanigans (might want a js framework as well)
ideally designed responsively from start
write apis for each main sport to grab fantasy league data
(will start with LoL, because dota has a fantasy_points api-call which is coming soon.
 Also I dislike the fantasyLCS system.
 Lol has no api so needs scraping, only would need to scrape once a day max so not an issue)
