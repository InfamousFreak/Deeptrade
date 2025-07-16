# backend/ml/fetch_stock.py
import yfinance as yf
import sys
import json

def get_stock_data(sym):
    # Download data
    data = yf.download(sym, period="7d", interval="1h")

    # Reset index
    data.reset_index(inplace=True)

    # Flatten MultiIndex columns: make col[0] if col is a tuple, else col
    data.columns = [col[0] if isinstance(col, tuple) else col for col in data.columns]

    # Convert datetime to string (yFinance uses 'Datetime' for intraday, 'Date' for daily)
    if "Datetime" in data.columns:
        data["Datetime"] = data["Datetime"].astype(str)
        datetime_col = "Datetime"
    elif "Date" in data.columns:
        data["Date"] = data["Date"].astype(str)
        datetime_col = "Date"
    else:
        raise Exception("No datetime column found")

    # Select only needed columns
    required = [datetime_col, "Open", "High", "Low", "Close"]
    missing = [col for col in required if col not in data.columns]
    if missing:
        raise Exception(f"Missing columns: {missing}")

    # Return clean JSON
    return data[required].to_dict(orient="records")

if __name__ == "__main__":
    symbol = sys.argv[1]
    try:
        stock_data = get_stock_data(symbol)
        print(json.dumps(stock_data))
    except Exception as e:
        print("❌ Error:", str(e))
        print(json.dumps({"error": str(e)}))
