import pandas as pd
from sqlalchemy import create_engine, text
import os

# ==========================================
# 1. 資料庫連線設定 (請修改這裡)
# ==========================================
DB_CONFIG = {

    "user": "postgres",
    "password": "123456",  # 您的密碼
    "host": "localhost",
    "port": "5432",
    "dbname": "ddnetone"     # 根據您的提示符 ddnetone=#，推測資料庫名稱為 ddnetone
}


EXCEL_FILE = '神秘活動新版.xlsx'

# ==========================================
# 2. Excel 解析邏輯 (針對您的檔案結構)
# ==========================================
def parse_excel():
    if not os.path.exists(EXCEL_FILE):
        print(f"❌ 找不到檔案 {EXCEL_FILE}")
        return pd.DataFrame()

    print("📖 正在解析 Excel 檔案...")
    all_data = []
    
    # 定義工作表與區塊
    sheets = ['NoviceModerateBrutalInsane', 'SoloDummyOldschoolRace', 'DDMAX EazyNextProNut', '多人圖Events']
    blocks = [(0, 5), (7, 12), (14, 19), (21, 26)]

    for sheet in sheets:
        try:
            df = pd.read_excel(EXCEL_FILE, sheet_name=sheet, header=None)
        except:
            continue

        for start, end in blocks:
            if start >= df.shape[1]: continue
            block = df.iloc[:, start:end+1].copy()
            
            # 1. 抓取難度 (Row 3, Index 3)
            # 邏輯：檢查第3列，若為空則檢查第2列

            diff = None
            if len(block) > 3:

                val = block.iloc[3, 0]
                if pd.notna(val) and str(val).strip() not in ['nan', 'Points', '完成地圖數']:
                    diff = str(val).strip()
            if not diff and len(block) > 2:

                val = block.iloc[2, 0]
                if pd.notna(val) and str(val).strip() not in ['nan', 'Points', '完成地圖數']:
                    diff = str(val).strip()
            
            if not diff or "總分" in diff: continue

            # 2. 抓取資料 (Map 欄位之後)
            start_row = 5
            for r in range(min(15, len(block))):
                if str(block.iloc[r, 0]).lower() == 'map':
                    start_row = r + 1
                    break

            # 3. 整理 DataFrame
            # Excel 欄位對應: Map(0), Pass(1), Save(2), Runner(3), Score(4), Note(5)
            data = block.iloc[start_row:].copy()
            # 重新命名為 DataFrame 暫時欄位
            data.columns = ['map_name', 'password', 'save_name', 'runner', 'score', 'note']
            data['difficulty'] = diff

            # 4. 過濾有效資料
            valid = data[
                (data['map_name'].notna()) & 
                (data['runner'].notna()) & 
                (~data['runner'].astype(str).str.strip().isin(['-', 'nan']))
            ].copy()

            if not valid.empty:
                # 轉型與清洗
                valid['score'] = pd.to_numeric(valid['score'], errors='coerce').fillna(0).astype(int)
                valid['note'] = valid['note'].fillna('')
                
                # 選取對應資料庫的欄位
                # map_records: difficulty, map_name, runner, score, note
                all_data.append(valid[['difficulty', 'map_name', 'runner', 'score', 'note']])

    if all_data:
        return pd.concat(all_data, ignore_index=True)
    return pd.DataFrame()

# ==========================================
# 3. 資料庫寫入與同步
# ==========================================
def main():
    # 1. 取得整理好的資料

    df = parse_excel()
    if df.empty:

        print("⚠️ 沒有提取到任何資料，程式結束。")
        return

    print(f"📊 提取到 {len(df)} 筆紀錄，準備寫入資料庫...")

    # 2. 建立連線
    conn_str = f"postgresql://{DB_CONFIG['user']}:{DB_CONFIG['password']}@{DB_CONFIG['host']}:{DB_CONFIG['port']}/{DB_CONFIG['dbname']}"
    engine = create_engine(conn_str)

    try:
        with engine.begin() as conn:
            # A. 清空舊資料 (這是一次性遷移，建議先清空以避免重複，若不想清空請註解掉這行)
            print("🧹 清空舊的 map_records 資料...")
            conn.execute(text("TRUNCATE TABLE map_records RESTART IDENTITY CASCADE;"))
            
            # B. 寫入 map_records
            print("💾 正在寫入 map_records ...")
            # method='multi' 可以加速大量寫入
            df.to_sql('map_records', conn, if_exists='append', index=False, method='multi')
            print(f"✅ map_records 寫入完成！")

            # C. 自動更新 players 表 (根據 map_records 重新統計)
            print("🔄 正在更新 players 統計表...")
            # 先清空 players
            conn.execute(text("TRUNCATE TABLE players RESTART IDENTITY;"))
            # 重新聚合插入
            sql_update_players = text("""
                INSERT INTO players (name, role, score_contribution, map_count, contribution_rate)
                SELECT 

                    runner,
                    'Agent',  -- 預設角色，您之後可以手動修改
                    SUM(score),

                    COUNT(*),

                    0 -- 暫時設為 0，稍後計算
                FROM map_records
                GROUP BY runner;
            """)

            conn.execute(sql_update_players)
            
            # 更新 contribution_rate (個人分數 / 總分)
            conn.execute(text("""

                UPDATE players 

                SET contribution_rate = ROUND((score_contribution / (SELECT SUM(score) FROM map_records)) * 100, 2)
                WHERE (SELECT SUM(score) FROM map_records) > 0;
            """))

            # D. 自動更新 summaries 表
            print("🔄 正在更新 summaries 統計表...")
            conn.execute(text("TRUNCATE TABLE summaries RESTART IDENTITY;"))

            conn.execute(text("""
                INSERT INTO summaries (current_score, target_score, completed_maps, last_update)
                VALUES (
                    (SELECT COALESCE(SUM(score), 0) FROM map_records),
                    10000, -- 您的目標分數
                    (SELECT COUNT(*) FROM map_records),
                    NOW()
                );

            """))

        print("-" * 30)
        print("🚀 資料庫遷移大成功！")
        print("1. map_records: 已填入 Excel 詳細資料")

        print("2. players: 已根據紀錄重新計算排名與積分")
        print("3. summaries: 已更新總分與進度")

    except Exception as e:
        print(f"❌ 資料庫操作失敗: {e}")

if __name__ == "__main__":
    main()
