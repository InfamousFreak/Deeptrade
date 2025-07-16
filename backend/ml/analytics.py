import yfinance as yf
import pandas as pd
import sys
import json

def calculate_indicators(symbol):
    try:
        # Explicitly set auto_adjust to True to avoid FutureWarning
        df = yf.download(symbol, period="1mo", interval="1d", auto_adjust=True)
        
        if df.empty:
            return {"error": f"No data available for symbol {symbol}"}
            
        df.reset_index(inplace=True)
        df['Date'] = df['Date'].astype(str)
    
        # Calculate SMAs
        df['SMA_7'] = df['Close'].rolling(window=7).mean()
        df['SMA_14'] = df['Close'].rolling(window=14).mean()
        
        # Calculate RSI (Relative Strength Index)
        delta = df['Close'].diff()
        gain = delta.where(delta > 0, 0)
        loss = -delta.where(delta < 0, 0)
        avg_gain = gain.rolling(window=14).mean()
        avg_loss = loss.rolling(window=14).mean()
        rs = avg_gain / (avg_loss + 1e-9)
        df['RSI'] = 100 - (100 / (1 + rs))
        
        # Drop rows with NaNs and get the latest row
        latest = df.dropna().iloc[-1]
        
        # Convert pandas Series values to native Python types using .item() for scalar conversion
        rsi = latest["RSI"].item()  # .item() is preferred for single-element Series
        sma_7 = latest["SMA_7"].item()
        sma_14 = latest["SMA_14"].item()
        date = latest["Date"]  # Already a string from astype(str) earlier
        
        # Generate signal
        signal = ""
        if rsi < 30:
            signal = "Buy (Oversold)"
        elif rsi > 70:
            signal = "Sell (Overbought)"
        else:
            signal = "Hold"
        
        return {
            "date": date,
            "sma_7": round(sma_7, 2),
            "sma_14": round(sma_14, 2),
            "rsi": round(rsi, 2),
            "signal": signal
        }
    except Exception as e:
        return {"error": str(e)}

if __name__ == "__main__":
    try:
        if len(sys.argv) < 2:
            print(json.dumps({"error": "Please provide a stock symbol"}))
            sys.exit(1)
            
        symbol = sys.argv[1]
        result = calculate_indicators(symbol)
        
        # Final check to ensure all values are JSON serializable
        for key in result:
            if isinstance(result[key], (pd.Series, pd.DataFrame)):
                if isinstance(result[key], pd.DataFrame):
                    result[key] = result[key].to_dict(orient='records')
                else:
                    result[key] = result[key].item() if len(result[key]) == 1 else result[key].tolist()
                    
        print(json.dumps(result))
            
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)
