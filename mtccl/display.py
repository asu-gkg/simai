import rich

def show_df(df):
    # 打印 df 的字段含义
    print("=== df 字段信息 ===")
    print(f"数据类型: {type(df)}")
    
    if hasattr(df, 'columns'):
        rich.print(f"\n=== DataFrame 列名及含义 ===")
        for i, col in enumerate(df.columns):
            rich.print(f"{i+1}. {col}")
        
        rich.print(f"\n=== 每层数据 ===")
        for i in range(len(df)):
            rich.print(f"\n--- 第 {i+1} 行 ---")
            for col in df.columns:
                value = df.iloc[i][col]
                rich.print(f"  {col}: {value}")
        rich.print(f"共 {len(df)} 行")
        
        rich.print(f"\n=== 列数据类型 ===")
        for col in df.columns:
            rich.print(f"  {col}: {df[col].dtype}")