IF NOT EXISTS (SELECT name
FROM sys.databases
WHERE name = 'BIS')
BEGIN
	CREATE DATABASE BIS;
END;
GO

USE [BIS]
GO

/****** Drop table contraints ******/
ALTER TABLE [dbo].[Bank] DROP CONSTRAINT [FK_Bank_Country]
GO
ALTER TABLE [dbo].[BankClient] DROP CONSTRAINT [FK_BankClient_Bank]
GO
ALTER TABLE [dbo].[BankEmployee] DROP CONSTRAINT [FK_BankEmployee_Bank]
GO
ALTER TABLE [dbo].[Transaction] DROP CONSTRAINT [FK_Transaction_Bank]
GO
ALTER TABLE [dbo].[Transaction] DROP CONSTRAINT [FK_Transaction_Bank1]
GO
ALTER TABLE [dbo].[Transaction] DROP CONSTRAINT [FK_Transaction_BankClient]
GO
ALTER TABLE [dbo].[Transaction] DROP CONSTRAINT [FK_Transaction_BankClient1]
GO
ALTER TABLE [dbo].[Transaction] DROP CONSTRAINT [FK_Transaction_TransactionType]
GO
ALTER TABLE [dbo].[TransactionPolicyStatus] DROP CONSTRAINT [FK_TransactionPolicyStatus_Policy]
GO
ALTER TABLE [dbo].[TransactionPolicyStatus] DROP CONSTRAINT [FK_TransactionPolicyStatus_Transaction]
GO
ALTER TABLE [dbo].[TransactionHistory] DROP CONSTRAINT [FK_TransactionHistory_Status]
GO
ALTER TABLE [dbo].[TransactionHistory] DROP CONSTRAINT [FK_TransactionHistory_Transaction]
GO
ALTER TABLE [dbo].[TransactionTypePolicy] DROP CONSTRAINT [FK_TransactionTypePolicy_TransactionType]
GO
ALTER TABLE [dbo].[TransactionTypePolicy] DROP CONSTRAINT [FK_TransactionTypePolicy_Policy]
GO
ALTER TABLE [dbo].[TransactionTypePolicy] DROP CONSTRAINT [FK_TransactionTypePolicy_Country]
GO
ALTER TABLE [dbo].[TransactionProof] DROP CONSTRAINT [FK_TransactionProof_Transaction]
GO

/****** Drop tables ******/
DROP TABLE [dbo].[Bank]
GO
DROP TABLE [dbo].[BankClient]
GO
DROP TABLE [dbo].[BankEmployee]
GO
DROP TABLE [dbo].[Policy]
GO
DROP TABLE [dbo].[Status]
GO
DROP TABLE [dbo].[Transaction]
GO
DROP TABLE [dbo].[TransactionHistory]
GO
DROP TABLE [dbo].[TransactionPolicyStatus]
GO
DROP TABLE [dbo].[TransactionProof]
GO
DROP TABLE [dbo].[TransactionType]
GO
DROP TABLE [dbo].[TransactionTypePolicy]
GO
DROP TABLE [dbo].[Country]
GO

/****** Object: Table [dbo].[Bank]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Bank]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[GlobalIdentifier] [nvarchar](512) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[Address] [nvarchar](500) NOT NULL,
	[CountryId] [int] NOT NULL,
	CONSTRAINT [PK_Bank] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[BankClient]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BankClient]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[GlobalIdentifier] [nvarchar](512) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[Address] [nvarchar](500) NOT NULL,
	[BankId] [bigint] NOT NULL,
	CONSTRAINT [PK_BankClient] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[BankEmployee]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BankEmployee]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[Username] [nvarchar](500) NOT NULL,
	[Password] [nvarchar](250) NOT NULL,
	[BankId] [bigint] NOT NULL,
	CONSTRAINT [PK_BankEmployee] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[Country]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Country]
(
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[CountryCode] [nvarchar](50) NULL,
	CONSTRAINT [PK_Country] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[Policy]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Policy]
(
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[CountryId] [int] NOT NULL,
	[Code] [nvarchar](1000) NOT NULL,
	[Name] [nvarchar](1000) NOT NULL,
	CONSTRAINT [PK_Policy] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[Status]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Status]
(
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](250) NOT NULL,
	CONSTRAINT [PK_Status] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[Transaction]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Transaction]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[OriginatorBank] [bigint] NOT NULL,
	[BeneficiaryBank] [bigint] NOT NULL,
	[Sender] [bigint] NOT NULL,
	[Receiver] [bigint] NOT NULL,
	[Currency] [nvarchar](50) NOT NULL,
	[Amount] [int] NOT NULL,
	[TransactionTypeId] [int] NOT NULL,
	[LoanId] [int] NOT NULL,
	CONSTRAINT [PK_Transaction] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[TransactionHistory]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionHistory]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[TransactionId] [bigint] NOT NULL,
	[StatusId] [int] NOT NULL,
	[Date] [datetime2](7) NOT NULL,
	CONSTRAINT [PK_TransactionHistory] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[TransactionPolicyStatus]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionPolicyStatus]
(
	[TransactionId] [bigint] NOT NULL,
	[PolicyId] [int] NOT NULL,
	[Status] [int] NOT NULL,
	CONSTRAINT [PK_TransactionPolicyStatus] PRIMARY KEY CLUSTERED 
(
	[TransactionId] ASC,
	[PolicyId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[TransactionProof]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionProof]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[TransactionId] [bigint] NOT NULL,
	[Proof] [nvarchar](500) NOT NULL,
	CONSTRAINT [PK_TransactionProof] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[TransactionType]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionType]
(
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Code] [nvarchar](100) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	CONSTRAINT [PK_Type] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object: Table [dbo].[TransactionTypePolicy]  Script Date: 25.10.2023. 13:17:44 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionTypePolicy]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[TransactionTypeId] [int] NOT NULL,
	[CountryId] [int] NOT NULL,
	[PolicyId] [int] NOT NULL,
	[Amount] [bigint] NOT NULL,
	[Checklist] [nvarchar](500) NOT NULL,
	CONSTRAINT [PK_TransactionTypePolicy] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Bank] WITH CHECK ADD CONSTRAINT [FK_Bank_Country] FOREIGN KEY([CountryId])
REFERENCES [dbo].[Country] ([Id])
GO
ALTER TABLE [dbo].[Bank] CHECK CONSTRAINT [FK_Bank_Country]
GO
ALTER TABLE [dbo].[BankClient] WITH CHECK ADD CONSTRAINT [FK_BankClient_Bank] FOREIGN KEY([BankId])
REFERENCES [dbo].[Bank] ([Id])
GO
ALTER TABLE [dbo].[BankClient] CHECK CONSTRAINT [FK_BankClient_Bank]
GO
ALTER TABLE [dbo].[BankEmployee] WITH CHECK ADD CONSTRAINT [FK_BankEmployee_Bank] FOREIGN KEY([BankId])
REFERENCES [dbo].[Bank] ([Id])
GO
ALTER TABLE [dbo].[BankEmployee] CHECK CONSTRAINT [FK_BankEmployee_Bank]
GO
ALTER TABLE [dbo].[Policy] WITH CHECK ADD CONSTRAINT [FK_Policy_Country] FOREIGN KEY([CountryId])
REFERENCES [dbo].[Country] ([Id])
GO
ALTER TABLE [dbo].[Policy] CHECK CONSTRAINT [FK_Policy_Country]
GO
ALTER TABLE [dbo].[Transaction] WITH CHECK ADD CONSTRAINT [FK_Transaction_Bank] FOREIGN KEY([OriginatorBank])
REFERENCES [dbo].[Bank] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_Bank]
GO
ALTER TABLE [dbo].[Transaction] WITH CHECK ADD CONSTRAINT [FK_Transaction_Bank1] FOREIGN KEY([BeneficiaryBank])
REFERENCES [dbo].[Bank] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_Bank1]
GO
ALTER TABLE [dbo].[Transaction] WITH CHECK ADD CONSTRAINT [FK_Transaction_BankClient] FOREIGN KEY([Sender])
REFERENCES [dbo].[BankClient] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_BankClient]
GO
ALTER TABLE [dbo].[Transaction] WITH CHECK ADD CONSTRAINT [FK_Transaction_BankClient1] FOREIGN KEY([Receiver])
REFERENCES [dbo].[BankClient] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_BankClient1]
GO
ALTER TABLE [dbo].[Transaction] WITH CHECK ADD CONSTRAINT [FK_Transaction_TransactionType] FOREIGN KEY([TransactionTypeId])
REFERENCES [dbo].[TransactionType] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_TransactionType]
GO
ALTER TABLE [dbo].[TransactionHistory] WITH CHECK ADD CONSTRAINT [FK_TransactionHistory_Status] FOREIGN KEY([StatusId])
REFERENCES [dbo].[Status] ([Id])
GO
ALTER TABLE [dbo].[TransactionHistory] CHECK CONSTRAINT [FK_TransactionHistory_Status]
GO
ALTER TABLE [dbo].[TransactionHistory] WITH CHECK ADD CONSTRAINT [FK_TransactionHistory_Transaction] FOREIGN KEY([TransactionId])
REFERENCES [dbo].[Transaction] ([Id])
GO
ALTER TABLE [dbo].[TransactionHistory] CHECK CONSTRAINT [FK_TransactionHistory_Transaction]
GO
ALTER TABLE [dbo].[TransactionPolicyStatus] WITH CHECK ADD CONSTRAINT [FK_TransactionPolicyStatus_Policy] FOREIGN KEY([PolicyId])
REFERENCES [dbo].[Policy] ([Id])
GO
ALTER TABLE [dbo].[TransactionPolicyStatus] CHECK CONSTRAINT [FK_TransactionPolicyStatus_Policy]
GO
ALTER TABLE [dbo].[TransactionPolicyStatus] WITH CHECK ADD CONSTRAINT [FK_TransactionPolicyStatus_Transaction] FOREIGN KEY([TransactionId])
REFERENCES [dbo].[Transaction] ([Id])
GO
ALTER TABLE [dbo].[TransactionPolicyStatus] CHECK CONSTRAINT [FK_TransactionPolicyStatus_Transaction]
GO
ALTER TABLE [dbo].[TransactionProof] WITH CHECK ADD CONSTRAINT [FK_TransactionProof_Transaction] FOREIGN KEY([TransactionId])
REFERENCES [dbo].[Transaction] ([Id])
GO
ALTER TABLE [dbo].[TransactionProof] CHECK CONSTRAINT [FK_TransactionProof_Transaction]
GO
ALTER TABLE [dbo].[TransactionTypePolicy] WITH CHECK ADD CONSTRAINT [FK_TransactionTypePolicy_Country] FOREIGN KEY([CountryId])
REFERENCES [dbo].[Country] ([Id])
GO
ALTER TABLE [dbo].[TransactionTypePolicy] CHECK CONSTRAINT [FK_TransactionTypePolicy_Country]
GO
ALTER TABLE [dbo].[TransactionTypePolicy] WITH CHECK ADD CONSTRAINT [FK_TransactionTypePolicy_Policy] FOREIGN KEY([PolicyId])
REFERENCES [dbo].[Policy] ([Id])
GO
ALTER TABLE [dbo].[TransactionTypePolicy] CHECK CONSTRAINT [FK_TransactionTypePolicy_Policy]
GO
ALTER TABLE [dbo].[TransactionTypePolicy] WITH CHECK ADD CONSTRAINT [FK_TransactionTypePolicy_TransactionType] FOREIGN KEY([TransactionTypeId])
REFERENCES [dbo].[TransactionType] ([Id])
GO
ALTER TABLE [dbo].[TransactionTypePolicy] CHECK CONSTRAINT [FK_TransactionTypePolicy_TransactionType]
GO