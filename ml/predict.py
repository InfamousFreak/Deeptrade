import yfinance as yf
import pandas as pd 
import numpy as np
import json
import sys
from xgboost import XGBClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score, f1_score

def get_dummy_sentiment(date):
    return np.random.uniform(-1, 1) 

def prepare_data(symbol):
    df = yf.download(symbol, period="6mo", interval="1d")
    df.reset_index(inplace=True)
    df = df.dropna()
    
    #add technical indicators
    df['SMA_7'] = df['Close'].rolling(window=7).mean()
    df['SMA_14'] = df['Close'].rolling(window=14).mean()
    delta = df['Close'].diff()
    gain = delta.where(delta > 0, 0)
    loss = -delta.where(delta < 0, 0)
    avg_gain = gain.rolling(window=14).mean()
    avg_loss = loss.rolling(window=14).mean()
    rs = avg_gain / (avg_loss + 1e-9)
    df['RSI'] = 100 - (100 / (1 + rs))
    
    df['Momentum'] = df['Close'] - df['Close'].shift(3)
    df['Sentiment'] = df['Date'].astype(str).apply(get_dummy_sentiment)
    
    df = df.dropna()
    
    #define features and labels 
    df['Target'] = (df['Close'].shift(-1) > df['Close']).astype(int)
    df = df.dropna()
    
    features = ['RSI', 'SMA_7', 'SMA_14', 'Momentum', 'Sentiment']
    X = df[features]
    y = df['Target']
    
    # Extract the date as a string - handle various formats and ensure it's a scalar value
    date_value = df.iloc[-1]['Date']
    
    # Convert to string using the appropriate method based on the data type
    if isinstance(date_value, pd.Series):
        # If it's a Series, extract the first value
        date_value = date_value.iloc[0]
    
    # Now convert to string based on type
    if pd.api.types.is_datetime64_any_dtype(date_value):
        # Ensure only the date part is displayed (no time component)
        latest_date = pd.Timestamp(date_value).date().strftime('%Y-%m-%d')
    else:
        latest_date = str(date_value)
        
    return X, y, df[features].iloc[-1:], latest_date #last row = today


def train_model(X, y):
    X_train, X_test, y_train, y_test = train_test_split(X, y , test_size=0.2, shuffle=False)
    model = XGBClassifier(use_label_encoder=False, eval_metric='logloss')
    model.fit(X_train, y_train)
    
    preds = model.predict(X_test)
    acc = accuracy_score(y_test, preds)
    f1 = f1_score(y_test, preds)
    
    return model, acc, f1

def predict_today(model, latest_data):
    pred = model.predict(latest_data)[0]
    conf = float(model.predict_proba(latest_data)[0][pred] * 100)
    return "UP" if pred == 1 else "DOWN", round(conf, 2)

if __name__ == "__main__":
    symbol = sys.argv[1]
    
    X, y, latest_data, latest_date = prepare_data(symbol)
    model, acc, f1 = train_model(X, y)
    direction, confidence = predict_today(model, latest_data)
    
    result = {
        "symbol": symbol.upper(),
        "date": latest_date,
        "prediction": direction,
        "confidence": float(confidence),
        "accuracy": float(round(float(acc) * 100, 2)),
        "f1_score": float(round(float(f1), 2)),
        "features_used": ["RSI", "SMA_7", "SMA_14", "Momentum", "Sentiment"]
    }
    
    print(json.dumps(result))
