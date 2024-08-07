DROP DATABASE BIS3
GO

IF NOT EXISTS (SELECT name
FROM sys.databases
WHERE name = 'BIS3')
BEGIN
	CREATE DATABASE BIS3;
END;
GO

USE [BIS3]
GO

/****** Object:  Table [dbo].[Bank]    Script Date: 03.06.2024. 14:06:14 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Bank]
(
	[GlobalIdentifier] [nvarchar](20) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[Address] [nvarchar](500) NOT NULL,
	[JurisdictionId] [nvarchar](6) NOT NULL,
	[BankTypeId] [int] NOT NULL,
	CONSTRAINT [PK_Bank] PRIMARY KEY CLUSTERED 
(
	[GlobalIdentifier] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[BankClient]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BankClient]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[GlobalIdentifier] [nvarchar](20) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[Address] [nvarchar](500) NOT NULL,
	[BankId] [nvarchar](20) NOT NULL,
	CONSTRAINT [PK_BankClient] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[BankEmployee]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BankEmployee]
(
	[Id] [nvarchar](20) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	[Username] [nvarchar](500) NOT NULL,
	[Password] [nvarchar](250) NOT NULL,
	[BankId] [nvarchar](20) NOT NULL,
	CONSTRAINT [PK_BankEmployee] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[BankType]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BankType]
(
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](250) NOT NULL,
	CONSTRAINT [PK_BankType] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Jurisdiction]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Jurisdiction]
(
	[Id] [nvarchar](6) NOT NULL,
	[Name] [nvarchar](500) NOT NULL,
	CONSTRAINT [PK_Jurisdiction] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Policy]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Policy]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[PolicyTypeId] [int] NOT NULL,
	[TransactionTypeId] [int] NOT NULL,
	[Owner] [nvarchar](20) NOT NULL,
	[PolicyEnforcingJurisdictionId] [nvarchar](6) NOT NULL,
	[OriginatingJurisdictionId] [nvarchar](6) NOT NULL,
	[BeneficiaryJurisdictionId] [nvarchar](6) NOT NULL,
	[Parameters] [nvarchar](max) NOT NULL,
	[IsPrivate] [bit] NOT NULL,
	[Latest] [bit] NOT NULL,
	CONSTRAINT [PK_TransactionTypePolicy] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[PolicyType]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[PolicyType]
(
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Code] [nvarchar](1000) NOT NULL,
	[Name] [nvarchar](1000) NOT NULL,
	CONSTRAINT [PK_Policy] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Status]    Script Date: 03.06.2024. 14:06:15 ******/
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
/****** Object:  Table [dbo].[Transaction]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Transaction]
(
	[Id] [nvarchar](64) NOT NULL,
	[OriginatorBankId] [nvarchar](20) NOT NULL,
	[BeneficiaryBankId] [nvarchar](20) NOT NULL,
	[SenderId] [bigint] NOT NULL,
	[ReceiverId] [bigint] NOT NULL,
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
/****** Object:  Table [dbo].[TransactionHistory]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionHistory]
(
	[TransactionId] [nvarchar](64) NOT NULL,
	[StatusId] [int] NOT NULL,
	[Date] [datetime2](7) NOT NULL,
	CONSTRAINT [PK_TransactionHistory_1] PRIMARY KEY CLUSTERED 
(
	[TransactionId] ASC,
	[StatusId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[TransactionPolicy]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionPolicy]
(
	[TransactionId] [nvarchar](64) NOT NULL,
	[PolicyId] [bigint] NOT NULL,
	[Status] [int] NOT NULL,
	[AdditionalParameters] [nvarchar](64),
	[Description] [nvarchar](200),
	CONSTRAINT [PK_TransactionPolicyStatus] PRIMARY KEY CLUSTERED 
(
	[TransactionId] ASC,
	[PolicyId] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[TransactionProof]    Script Date: 03.06.2024. 14:06:15 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TransactionProof]
(
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[TransactionId] [nvarchar](64) NOT NULL,
	[Proof] [nvarchar](500) NOT NULL,
	CONSTRAINT [PK_TransactionProof] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[TransactionType]    Script Date: 03.06.2024. 14:06:15 ******/
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
ALTER TABLE [dbo].[Bank]  WITH CHECK ADD  CONSTRAINT [FK_Bank_BankType] FOREIGN KEY([BankTypeId])
REFERENCES [dbo].[BankType] ([Id])
GO
ALTER TABLE [dbo].[Bank] CHECK CONSTRAINT [FK_Bank_BankType]
GO
ALTER TABLE [dbo].[Bank]  WITH CHECK ADD  CONSTRAINT [FK_Bank_Jurisdiction] FOREIGN KEY([JurisdictionId])
REFERENCES [dbo].[Jurisdiction] ([Id])
GO
ALTER TABLE [dbo].[Bank] CHECK CONSTRAINT [FK_Bank_Jurisdiction]
GO
ALTER TABLE [dbo].[BankClient]  WITH CHECK ADD  CONSTRAINT [FK_BankClient_Bank] FOREIGN KEY([BankId])
REFERENCES [dbo].[Bank] ([GlobalIdentifier])
GO
ALTER TABLE [dbo].[BankClient] CHECK CONSTRAINT [FK_BankClient_Bank]
GO
ALTER TABLE [dbo].[BankEmployee]  WITH CHECK ADD  CONSTRAINT [FK_BankEmployee_Bank] FOREIGN KEY([BankId])
REFERENCES [dbo].[Bank] ([GlobalIdentifier])
GO
ALTER TABLE [dbo].[BankEmployee] CHECK CONSTRAINT [FK_BankEmployee_Bank]
GO
ALTER TABLE [dbo].[Policy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionTypePolicy_Jurisdiction] FOREIGN KEY([PolicyEnforcingJurisdictionId])
REFERENCES [dbo].[Jurisdiction] ([Id])
GO
ALTER TABLE [dbo].[Policy] CHECK CONSTRAINT [FK_TransactionTypePolicy_Jurisdiction]
GO
ALTER TABLE [dbo].[Policy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionTypePolicy_Jurisdiction1] FOREIGN KEY([OriginatingJurisdictionId])
REFERENCES [dbo].[Jurisdiction] ([Id])
GO
ALTER TABLE [dbo].[Policy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionTypePolicy_Jurisdiction2] FOREIGN KEY([BeneficiaryJurisdictionId])
REFERENCES [dbo].[Jurisdiction] ([Id])
GO
ALTER TABLE [dbo].[Policy] CHECK CONSTRAINT [FK_TransactionTypePolicy_Jurisdiction1]
GO
ALTER TABLE [dbo].[Policy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionTypePolicy_Policy] FOREIGN KEY([PolicyTypeId])
REFERENCES [dbo].[PolicyType] ([Id])
GO
ALTER TABLE [dbo].[Policy] CHECK CONSTRAINT [FK_TransactionTypePolicy_Policy]
GO
ALTER TABLE [dbo].[Policy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionTypePolicy_TransactionType] FOREIGN KEY([TransactionTypeId])
REFERENCES [dbo].[TransactionType] ([Id])
GO
ALTER TABLE [dbo].[Policy]  WITH CHECK ADD  CONSTRAINT [FK_Owner_Bank] FOREIGN KEY([Owner])
REFERENCES [dbo].[Bank] ([GlobalIdentifier])
GO
ALTER TABLE [dbo].[Policy] CHECK CONSTRAINT [FK_TransactionTypePolicy_TransactionType]
GO
ALTER TABLE [dbo].[Transaction]  WITH CHECK ADD  CONSTRAINT [FK_Transaction_Bank] FOREIGN KEY([OriginatorBankId])
REFERENCES [dbo].[Bank] ([GlobalIdentifier])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_Bank]
GO
ALTER TABLE [dbo].[Transaction]  WITH CHECK ADD  CONSTRAINT [FK_Transaction_Bank1] FOREIGN KEY([BeneficiaryBankId])
REFERENCES [dbo].[Bank] ([GlobalIdentifier])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_Bank1]
GO
ALTER TABLE [dbo].[Transaction]  WITH CHECK ADD  CONSTRAINT [FK_Transaction_BankClient] FOREIGN KEY([SenderId])
REFERENCES [dbo].[BankClient] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_BankClient]
GO
ALTER TABLE [dbo].[Transaction]  WITH CHECK ADD  CONSTRAINT [FK_Transaction_BankClient1] FOREIGN KEY([ReceiverId])
REFERENCES [dbo].[BankClient] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_BankClient1]
GO
ALTER TABLE [dbo].[Transaction]  WITH CHECK ADD  CONSTRAINT [FK_Transaction_TransactionType] FOREIGN KEY([TransactionTypeId])
REFERENCES [dbo].[TransactionType] ([Id])
GO
ALTER TABLE [dbo].[Transaction] CHECK CONSTRAINT [FK_Transaction_TransactionType]
GO
ALTER TABLE [dbo].[TransactionHistory]  WITH CHECK ADD  CONSTRAINT [FK_TransactionHistory_Status] FOREIGN KEY([StatusId])
REFERENCES [dbo].[Status] ([Id])
GO
ALTER TABLE [dbo].[TransactionHistory] CHECK CONSTRAINT [FK_TransactionHistory_Status]
GO
ALTER TABLE [dbo].[TransactionHistory]  WITH CHECK ADD  CONSTRAINT [FK_TransactionHistory_Transaction] FOREIGN KEY([TransactionId])
REFERENCES [dbo].[Transaction] ([Id])
GO
ALTER TABLE [dbo].[TransactionHistory] CHECK CONSTRAINT [FK_TransactionHistory_Transaction]
GO
ALTER TABLE [dbo].[TransactionPolicy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionPolicy_Policy] FOREIGN KEY([PolicyId])
REFERENCES [dbo].[Policy] ([Id])
GO
ALTER TABLE [dbo].[TransactionPolicy] CHECK CONSTRAINT [FK_TransactionPolicy_Policy]
GO
ALTER TABLE [dbo].[TransactionPolicy]  WITH CHECK ADD  CONSTRAINT [FK_TransactionPolicyStatus_Transaction] FOREIGN KEY([TransactionId])
REFERENCES [dbo].[Transaction] ([Id])
GO
ALTER TABLE [dbo].[TransactionPolicy] CHECK CONSTRAINT [FK_TransactionPolicyStatus_Transaction]
GO
ALTER TABLE [dbo].[TransactionProof]  WITH CHECK ADD  CONSTRAINT [FK_TransactionProof_Transaction] FOREIGN KEY([TransactionId])
REFERENCES [dbo].[Transaction] ([Id])
GO
ALTER TABLE [dbo].[TransactionProof] CHECK CONSTRAINT [FK_TransactionProof_Transaction]
GO
