import psycopg2
import os

# ==========================================
# 1. 資料庫連線設定
# ==========================================
DB_CONFIG = {
    "host": "localhost",
    "user": "postgres",
    "password": "123456",  # 請修改您的密碼
    "dbname": "ddnetone",    # 資料庫名稱
    "port": "5432"
}

INPUT_FILE = 'maps.txt'


def calculate_score(star):
    """
    分數計算公式: 
    1星 -> 34
    2星 -> 38
    3星 -> 42
    4星 -> 46
    5星 -> 50
    公式: 30 + (星數 * 4)
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


        # 1. 取得資料庫現有地圖 (避免重複匯入)
        cursor.execute("SELECT map_name FROM map_records")
        existing_maps = set(row[0] for row in cursor.fetchall())
        print(f"📊 資料庫現有地圖數量: {len(existing_maps)}")

        # 2. 解析 maps.txt
        new_records = []
        
        with open(INPUT_FILE, 'r', encoding='utf-8') as f:
            lines = f.readlines()

        print("📖 正在解析並計算分數...")
        
        for line in lines:
            line = line.strip()
            if not line or "───" in line:
                continue


            parts = line.split('|')
            

            # 格式檢查: 必須是 "星數|地圖名|..."
            if len(parts) >= 2 and parts[0].isdigit():
                star = int(parts[0])
                map_name = parts[1].strip()
                
                # 計算分數
                map_score = calculate_score(star)

                # ★ 如果資料庫沒有這張圖，才加入
                if map_name not in existing_maps:
                    new_records.append({
                        "map_name": map_name,
                        "score": map_score
                    })
                    existing_maps.add(map_name)

        # 3. 匯入資料庫
        if not new_records:
            print("✅ 沒有發現新地圖，資料庫已是最新。")
        else:
            print(f"🚀 發現 {len(new_records)} 張新地圖，準備匯入...")
            
            # SQL: 直接將計算好的分數寫入 score 欄位，Status 設為 0
            insert_query = """
            INSERT INTO map_records 
            (difficulty, map_name, runner, score, note, status) 
            VALUES ('INSANE', %s, '', %s, '', 0)
            """
            
            data_to_insert = [
                (r['map_name'], r['score']) 
                for r in new_records
            ]

            cursor.executemany(insert_query, data_to_insert)
            conn.commit()

            
            print(f"🎉 成功匯入 {cursor.rowcount} 筆新地圖！")
            print("   (欄位設定: Status=0, Runner='', Note='', Score=地圖分數)")

    except Exception as e:
        print(f"❌ 發生錯誤: {e}")
    finally:
        if conn:

            cursor.close()
            conn.close()
            print("🔒 連線已關閉")


if __name__ == "__main__":
    main()
