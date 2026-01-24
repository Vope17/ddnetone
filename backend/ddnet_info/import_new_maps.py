import psycopg2
import os
from pathlib import Path
from dotenv import load_dotenv

# 取得目前檔案的父目錄的父目錄 (假設腳本在根目錄的子資料夾內)
# 或者直接使用 Path.cwd() 取得目前執行指令的工作目錄
env_path = Path(__file__).resolve().parent.parent.parent / '.env'

# 讀取指定的 .env 路徑
load_dotenv(dotenv_path=env_path)

try:
    host = os.environ['DB_HOST']
    user = os.environ['DB_USER']
    password = os.environ['DB_PASSWORD']
    dbname = os.environ['DB_NAME']
    port = os.environ['DB_PORT']
    # ... 其餘變數
    print(f"成功從 {env_path} 讀取設定")
except KeyError as e:
    print(f"找不到 DB 環境變數: {e}")

# ==========================================
# 1. 資料庫連線設定
# ==========================================
DB_CONFIG = {
    "host": host,
    "user": user,
    "password": password,  # ★ 請修改為您的密碼
    "dbname": dbname,    # 資料庫名稱
    "port": port
}

INPUT_FILE = './insane_maps.txt'
TARGET_DIFFICULTY = 'INSANE' # 統一難度

def calculate_points(star):
    """
    分數計算公式: 3 + (星數 * 2)
    """
    try:
        s = int(star)
        return 30 + (s * 4)
    except:
        return 0

def main():
    if not os.path.exists(INPUT_FILE):
        print(f"❌ 找不到檔案: {INPUT_FILE}")
        return

    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cursor = conn.cursor()
        print("🔌 資料庫連線成功")

        # ---------------------------------------------------------
        # 1. 先讀取 maps.txt 並整理成字典 { '地圖名': 計算後的分數 }
        # ---------------------------------------------------------
        file_map_data = {} # 格式: {'MapName': 42, ...}
        
        with open(INPUT_FILE, 'r', encoding='utf-8') as f:
            lines = f.readlines()

        print("📖 正在解析 maps.txt ...")
        for line in lines:
            line = line.strip()
            # 跳過空行或標題
            if not line or "───" in line:
                continue

            parts = line.split('|')
            # 格式: 星數|地圖名|...
            if len(parts) >= 2 and parts[0].isdigit():
                star = int(parts[0])
                map_name = parts[1].strip()
                points = calculate_points(star)
                
                file_map_data[map_name] = points

        print(f"📄 檔案中共有 {len(file_map_data)} 張地圖資料")

        # ---------------------------------------------------------
        # 2. 取得資料庫現況 (MapName 和 Points)
        # ---------------------------------------------------------
        cursor.execute("SELECT map_name, points FROM map_records")
        # 轉成字典: { 'MapName': 目前資料庫的points }
        db_map_data = {row[0]: row[1] for row in cursor.fetchall()}
        
        print(f"📊 資料庫現有 {len(db_map_data)} 張地圖")

        # ---------------------------------------------------------
        # 3. 分類：哪些要新增？哪些要更新？
        # ---------------------------------------------------------
        to_insert = []
        to_update = []

        for map_name, correct_points in file_map_data.items():
            if map_name not in db_map_data:
                # 情況 A: 資料庫沒有 -> 新增

                to_insert.append((map_name, correct_points))
            
            else:
                # ★★★ 修改處：只要資料庫有，不管 points 是多少，都加入更新清單 ★★★

                to_update.append((correct_points, map_name))

        # ---------------------------------------------------------

        # 4. 執行資料庫操作
        # ---------------------------------------------------------

        # A. 執行新增 (Insert)
        if to_insert:
            print(f"🚀 發現 {len(to_insert)} 張新地圖，正在新增...")
            insert_query = """
            INSERT INTO map_records 
            (difficulty, map_name, runner, points, score, note, status) 
            VALUES (%s, %s, '', %s, 0, '', 0)
            """
            # 構建參數: (Difficulty, MapName, Points)
            insert_params = [(TARGET_DIFFICULTY, name, pts) for name, pts in to_insert]
            cursor.executemany(insert_query, insert_params)
            print(f"   ✅ 已新增 {cursor.rowcount} 筆資料")
        else:
            print("✅ 沒有需要新增的地圖")

        # B. 執行更新 (Update)
        if to_update:
            print(f"🔧 發現 {len(to_update)} 張現有地圖，正在強制更新分數...")
            update_query = """
            UPDATE map_records 
            SET points = %s 
            WHERE map_name = %s
            """
            # 參數順序必須對應 SQL: (Points, MapName)
            cursor.executemany(update_query, to_update)
            print(f"   ✅ 已更新 {cursor.rowcount} 筆資料的 Points")
        else:
            print("✅ 沒有需要更新的地圖")

        # 提交變更
        conn.commit()

    except Exception as e:
        print(f"❌ 發生錯誤: {e}")
        if conn:
            conn.rollback() # 發生錯誤時回滾
    finally:
        if conn:
            cursor.close()
            conn.close()
            print("🔒 資料庫連線已關閉")

if __name__ == "__main__":
    main()
