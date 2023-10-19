USE [BIS]
GO

SET IDENTITY_INSERT [dbo].[Country] ON 
INSERT [dbo].[Country] ([Id], [Name], [CountryCode]) VALUES (1, N'Malaysia', N'MYS')
INSERT [dbo].[Country] ([Id], [Name], [CountryCode]) VALUES (2, N'Singapore', N'SGP')
INSERT [dbo].[Country] ([Id], [Name], [CountryCode]) VALUES (3, N'Switzerland', N'CHE')
INSERT [dbo].[Country] ([Id], [Name], [CountryCode]) VALUES (4, N'Germany', N'DEU')
SET IDENTITY_INSERT [dbo].[Country] OFF
GO

SET IDENTITY_INSERT [dbo].[Bank] ON
INSERT [dbo].[Bank] ([Id], [GlobalIdentifier], [Name], [Address], [CountryId]) VALUES (1, N'AAAA-BB-CC-123', N'BNM', N'Jalan Dato Onn 50480 Kuala Lumpur', 1)
INSERT [dbo].[Bank] ([Id], [GlobalIdentifier], [Name], [Address], [CountryId]) VALUES (2, N'EEEE-NN-MM-123', N'MAS', N'10 Shenton Wy', 2)
INSERT [dbo].[Bank] ([Id], [GlobalIdentifier], [Name], [Address], [CountryId]) VALUES (3, N'TTTT-PP-XX-123', N'SNB', N'Börsenstrasse 15, Zurich', 3)
INSERT [dbo].[Bank] ([Id], [GlobalIdentifier], [Name], [Address], [CountryId]) VALUES (4, N'QQQQ-LL-AA-123', N'DB', N'Wilhelm-Epstein-Strasse 14', 4)
SET IDENTITY_INSERT [dbo].[Bank] OFF
GO

SET IDENTITY_INSERT [dbo].[BankClient] ON 
INSERT [dbo].[BankClient] ([Id], [GlobalIdentifier], [Name], [Address], [BankId]) VALUES (1, N'AAAA-BB-CC-321', N'Lạc Thị Mai Châu', N'Jalan Dato Onn', 1)
INSERT [dbo].[BankClient] ([Id], [GlobalIdentifier], [Name], [Address], [BankId]) VALUES (2, N'EEEE-NN-MM-321', N'Jake Headlam', N'Jln. Bukit Ho Swee', 2)
INSERT [dbo].[BankClient] ([Id], [GlobalIdentifier], [Name], [Address], [BankId]) VALUES (3, N'TTTT-PP-XX-321', N'Bernd Braun', N'Hottingerstrasse 16', 3)
INSERT [dbo].[BankClient] ([Id], [GlobalIdentifier], [Name], [Address], [BankId]) VALUES (4, N'QQQQ-LL-AA-321', N'Jürgen Huber', N'Karlstraße 79', 4)
SET IDENTITY_INSERT [dbo].[BankClient] OFF
GO

SET IDENTITY_INSERT [dbo].[BankEmployee] ON 
INSERT [dbo].[BankEmployee] ([Id], [Name], [Username], [Password], [BankId]) VALUES (1, N'Văn Hải Thao', N'BNMadmin', N'admin', 1)
INSERT [dbo].[BankEmployee] ([Id], [Name], [Username], [Password], [BankId]) VALUES (2, N'Amy Walton', N'MASadmin', N'admin', 2)
INSERT [dbo].[BankEmployee] ([Id], [Name], [Username], [Password], [BankId]) VALUES (3, N'Petra Traugott', N'SNBadmin', N'admin', 3)
INSERT [dbo].[BankEmployee] ([Id], [Name], [Username], [Password], [BankId]) VALUES (4, N'Robert Brauer', N'DBadmin', N'admin', 4)
SET IDENTITY_INSERT [dbo].[BankEmployee] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionType] ON 
INSERT [dbo].[TransactionType] ([Id], [Name]) VALUES (1, N'31111')
INSERT [dbo].[TransactionType] ([Id], [Name]) VALUES (2, N'31121')
SET IDENTITY_INSERT [dbo].[TransactionType] OFF
GO

SET IDENTITY_INSERT [dbo].[Policy] ON 
INSERT [dbo].[Policy] ([Id], [Name]) VALUES (1, N'Capital Flow Management')
INSERT [dbo].[Policy] ([Id], [Name]) VALUES (2, N'Saction Check List')
SET IDENTITY_INSERT [dbo].[Policy] OFF
GO

INSERT [dbo].[TransactionTypePolicy] ([TransactionTypeId], [PolicyId], [CountryId], [Amount]) VALUES (1, 1, 1, 100000)
INSERT [dbo].[TransactionTypePolicy] ([TransactionTypeId], [PolicyId], [CountryId], [Amount]) VALUES (1, 2, 1, 1000000)
INSERT [dbo].[TransactionTypePolicy] ([TransactionTypeId], [PolicyId], [CountryId], [Amount]) VALUES (2, 1, 1, 100000)
INSERT [dbo].[TransactionTypePolicy] ([TransactionTypeId], [PolicyId], [CountryId], [Amount]) VALUES (2, 2, 1, 1000000)
GO

SET IDENTITY_INSERT [dbo].[Transaction] ON 
INSERT [dbo].[Transaction] ([Id], [OriginatorBank], [BeneficiaryBank], [Sender], [Receiver], [Currency], [Amount], [TransactionTypeId]) VALUES (1, 1, 2, 1, 2, N'MYR', 100000, 1)
INSERT [dbo].[Transaction] ([Id], [OriginatorBank], [BeneficiaryBank], [Sender], [Receiver], [Currency], [Amount], [TransactionTypeId]) VALUES (2, 3, 1, 3, 1, N'MYR', 500000, 2)
INSERT [dbo].[Transaction] ([Id], [OriginatorBank], [BeneficiaryBank], [Sender], [Receiver], [Currency], [Amount], [TransactionTypeId]) VALUES (3, 4, 1, 4, 1, N'USD', 1000000, 1)
SET IDENTITY_INSERT [dbo].[Transaction] OFF
GO

SET IDENTITY_INSERT [dbo].[Status] ON 
INSERT [dbo].[Status] ([Id], [Name]) VALUES (1, N'TransactionCreated')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (2, N'PoliciesApplied')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (3, N'ProofRequested')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (4, N'ProofReceived')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (5, N'ProofInvalid')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (6, N'AssetSent')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (7, N'TransactionCompleted')
INSERT [dbo].[Status] ([Id], [Name]) VALUES (8, N'TransactionCanceled')
SET IDENTITY_INSERT [dbo].[Status] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionHistory] ON 
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (1, 1, 1, CAST(N'2023-10-10T00:00:00.0000000' AS DateTime2))
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (2, 1, 2, CAST(N'2023-10-11T00:00:00.0000000' AS DateTime2))
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (3, 1, 3, CAST(N'2023-10-12T00:00:00.0000000' AS DateTime2))
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (4, 1, 4, CAST(N'2023-10-13T00:00:00.0000000' AS DateTime2))
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (5, 1, 7, CAST(N'2023-10-13T15:00:00.0000000' AS DateTime2))
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (6, 2, 1, CAST(N'2023-10-14T00:00:00.0000000' AS DateTime2))
INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (7, 3, 1, CAST(N'2023-10-15T00:00:00.0000000' AS DateTime2))
SET IDENTITY_INSERT [dbo].[TransactionHistory] OFF
GO
