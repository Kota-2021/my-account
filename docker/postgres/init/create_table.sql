-- カテゴリーマスタ
create table m_categories (
	category_id smallint PRIMARY KEY,
	category_name varchar(6) NOT NULL
);

-- 勘定科目マスタ
create table m_subjects (
	subject_code smallint PRIMARY KEY,
	subject_name varchar(20) NOT NULL
);

-- 帳票マスタ
create table m_books (  
	book_code smallint PRIMARY KEY,
	book_name varchar(6) NOT NULL
);

-- 出納帳データ
create table t_cashbook (
	cashbook_id SERIAL PRIMARY KEY,
	cashbook_date date NOT NULL,
	item varchar(20),
	withdrawal decimal NOT NULL,
	deposit decimal NOT NULL,
	balance decimal NOT NULL,
	remarks varchar(20),
	book_code smallint NOT NULL,
	book_year smallint NOT NULL,

	CONSTRAINT fk_t_cashbook_book
		FOREIGN KEY (book_code)
		REFERENCES m_books(book_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL
);

-- 仕訳帳データ
create table t_journal (
	journal_id SERIAL PRIMARY KEY,
	journal_date date NOT NULL,
	withdrawal decimal NOT NULL,
	deposit decimal NOT NULL,
	subject_code smallint NOT NULL,
	item varchar(20),
	customer varchar(20) NOT NULL,
	evidence varchar(6),
	memo varchar(20),
	book_code smallint NOT NULL,
	category_id smallint NOT NULL,
	fiscal_year smallint NOT NULL,

	CONSTRAINT fk_t_journal_subject
		FOREIGN KEY (subject_code)
		REFERENCES m_subjects(subject_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,

	CONSTRAINT fk_t_journal_book
		FOREIGN KEY (book_code)
		REFERENCES m_books(book_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,
	
	CONSTRAINT fk_t_journal_category
		FOREIGN KEY (category_id)
		REFERENCES m_categories(category_id)
		ON UPDATE CASCADE
		ON DELETE SET NULL
);

-- 未収金データ
create table t_receivable (
	receivable_id SERIAL PRIMARY KEY,
	receivable_date date NOT NULL,
	withdrawal decimal NOT NULL,
	deposit decimal NOT NULL,
	subject_code smallint NOT NULL,
	item varchar(20),
	customer varchar(20) NOT NULL,
	evidence varchar(6),
	memo varchar(20),
	book_code smallint NOT NULL,
	category_id smallint NOT NULL,
	fiscal_year smallint NOT NULL,

	CONSTRAINT fk_t_receivable_subject
		FOREIGN KEY (subject_code)
		REFERENCES m_subjects(subject_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,

	CONSTRAINT fk_t_receivable_book
		FOREIGN KEY (book_code)
		REFERENCES m_books(book_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,

	CONSTRAINT fk_t_receivable_category
		FOREIGN KEY (category_id)
		REFERENCES m_categories(category_id)
		ON UPDATE CASCADE
		ON DELETE SET NULL
);

-- 未払金データ
create table t_payable (
	payable_id SERIAL PRIMARY KEY,
	payable_date date NOT NULL,
	withdrawal decimal NOT NULL,
	deposit decimal NOT NULL,
	subject_code smallint NOT NULL,
	item varchar(20),
	customer varchar(20) NOT NULL,
	evidence varchar(6),
	memo varchar(20),
	book_code smallint NOT NULL,
	category_id smallint NOT NULL,
	fiscal_year smallint NOT NULL,

	CONSTRAINT fk_t_payable_subject
		FOREIGN KEY (subject_code)
		REFERENCES m_subjects(subject_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,

	CONSTRAINT fk_t_payable_book
		FOREIGN KEY (book_code)
		REFERENCES m_books(book_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,

	CONSTRAINT fk_t_payable_category
		FOREIGN KEY (category_id)
		REFERENCES m_categories(category_id)
		ON UPDATE CASCADE
		ON DELETE SET NULL
);

-- 予算・決算データ
create table t_buget_financial_data (
	buget_financial_data_id SERIAL PRIMARY KEY,
	subject_code smallint NOT NULL,
	budget decimal NOT NULL,
	result decimal NOT NULL,
	difference decimal NOT NULL,
	category_id smallint NOT NULL,
	buget_fiscal_year smallint NOT NULL,
	UNIQUE (subject_code, category_id, buget_fiscal_year),

	CONSTRAINT fk_t_buget_financial_data_subject
		FOREIGN KEY (subject_code)
		REFERENCES m_subjects(subject_code)
		ON UPDATE CASCADE
		ON DELETE SET NULL,

	CONSTRAINT fk_t_buget_financial_data_category
		FOREIGN KEY (category_id)
		REFERENCES m_categories(category_id)
		ON UPDATE CASCADE
		ON DELETE SET NULL
);
