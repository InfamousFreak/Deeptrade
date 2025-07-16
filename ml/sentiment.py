from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import sys, json

analyzer = SentimentIntensityAnalyzer()
headlines = json.loads(sys.argv[1])
results = []

for h in headlines:
    score = analyzer.polarity_scores(h)["compound"]  # Range: -1 to 1
    results.append({"headline": h, "score": score})

print(json.dumps(results))
