import yfinance as yf
import pandas as pd
import numpy as np
import sys
import json

def backtest_strategy(symbol):
    try:
        # Use Ticker class instead of download
        symbol = str(symbol).strip()
        print(json.dumps({"debug": f"Downloading data for ticker: {symbol}"}))
        
        # Create a Ticker object and fetch the data
        ticker = yf.Ticker(symbol)
        df = ticker.history(period="1mo", interval="1d", auto_adjust=True)
        
        if df.empty:
            return {
                "error": f"No data available for symbol {symbol}"
            }
        
        df = df.copy()
        df.reset_index(inplace=True)
        df['Date'] = df['Date'].astype(str)
        
        df['Close'] = pd.to_numeric(df['Close'], errors='coerce')
        
        # Calculate Simple Moving Averages
        df['SMA_7'] = df['Close'].rolling(window=7).mean()
        df['SMA_14'] = df['Close'].rolling(window=14).mean()
        
        # Calculate RSI - Fixed version
        delta = df['Close'].diff()
        gain = delta.where(delta > 0, 0)
        loss = -delta.where(delta < 0, 0)
        avg_gain = gain.rolling(window=14).mean()
        avg_loss = loss.rolling(window=14).mean()
        
        # Avoid division by zero and handle the calculation properly
        avg_loss_safe = avg_loss.replace(0, np.nan)
        rs = avg_gain / avg_loss_safe
        df['RSI'] = 100 - (100 / (1 + rs))
        
        # Fill any NaN values in RSI with 50 (neutral value)
        df['RSI'] = df['RSI'].fillna(50)
        
        df['SMA_7'] = pd.to_numeric(df['SMA_7'], errors='coerce')
        df['SMA_14'] = pd.to_numeric(df['SMA_14'], errors='coerce')
        df['RSI'] = pd.to_numeric(df['RSI'], errors='coerce')
        
        df = df.dropna(subset=['Close', 'SMA_7', 'SMA_14'])
        
        capital = 10000
        shares = 0
        cash = capital
        trades = []
        
        for i in range(1, len(df)):
            
            rsi_val = df.iloc[i]["RSI"]
            sma_7_val = df.iloc[i]["SMA_7"] 
            sma_14_val = df.iloc[i]["SMA_14"]
            close_price = df.iloc[i]["Close"]
            date = df.iloc[i]["Date"]
            
            # Convert numpy values to Python scalars
            if hasattr(rsi_val, 'item'):
                rsi_val = rsi_val.item()
            if hasattr(sma_7_val, 'item'):
                sma_7_val = sma_7_val.item()
            if hasattr(sma_14_val, 'item'):
                sma_14_val = sma_14_val.item()
            if hasattr(close_price, 'item'):
                close_price = close_price.item()
            
            # Check for NaN values
            if pd.isna(rsi_val) or pd.isna(sma_7_val) or pd.isna(sma_14_val) or pd.isna(close_price):
                continue
            
            # Trading signals
            buy_signal = (rsi_val < 30) and (sma_7_val > sma_14_val)
            sell_signal = (rsi_val > 70) and (sma_7_val < sma_14_val)
            
            # Execute trades
            if buy_signal and shares == 0 and cash > 0:
                shares_to_buy = int(cash // close_price)
                if shares_to_buy > 0:
                    shares += shares_to_buy
                    cash -= shares_to_buy * close_price
                    trades.append({
                        "date": date,
                        "action": "BUY",
                        "shares": shares_to_buy,
                        "price": round(close_price, 2),
                        "total": round(shares_to_buy * close_price, 2)
                    })
                    
            elif sell_signal and shares > 0:
                # Sell signal - sell all shares
                cash += shares * close_price
                trades.append({
                    "date": date,
                    "action": "SELL",
                    "shares": shares,
                    "price": round(close_price, 2),
                    "total": round(shares * close_price, 2)
                })
                shares = 0
                
        # Calculate final portfolio value
        if len(df) > 0:
            final_price = df.iloc[-1]["Close"]
            if hasattr(final_price, 'item'):
                final_price = final_price.item()
            final_portfolio_value = cash + (shares * final_price)
        else:
            final_portfolio_value = cash
            
        # Calculate returns
        total_return = final_portfolio_value - capital
        return_percent = (total_return / capital) * 100
        
        result = {
            "symbol": symbol,
            "initial_capital": capital,
            "final_portfolio_value": round(final_portfolio_value, 2),
            "total_return": round(total_return, 2),
            "return_percent": round(return_percent, 2),
            "total_trades": len(trades),
            "final_cash": round(cash, 2),
            "final_shares": shares,
            "trades": trades[-10:] if len(trades) > 10 else trades 
        }  
        
        return result
    
    except Exception as e:
        return {"error": str(e)}


if __name__ == "__main__":
    try:
        if len(sys.argv) < 2:
            print(json.dumps({"error": "Please provide a stock symbol"}))
            sys.exit(1)
        
        symbol = sys.argv[1]
        result = backtest_strategy(symbol)
        
        # Handle pandas objects in result
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