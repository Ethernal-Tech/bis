USE [BIS2]
GO
SET IDENTITY_INSERT [dbo].[BankType] ON

INSERT [dbo].[BankType]
    ([Id], [Name])
VALUES
    (1, N'Commercial')
INSERT [dbo].[BankType]
    ([Id], [Name])
VALUES
    (2, N'Central')
SET IDENTITY_INSERT [dbo].[BankType] OFF
GO

INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'SG', N'Singapore')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'MY', N'Malaysia')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'AU', N'Australia')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'KR', N'South Korea')

INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'54930035WQZLGC45RZ35', N'Monetary Authority of Singapore', N'10 Shenton Wy', N'SG', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'5493003GYPR7VI37GG77', N'Reserve Bank of Australia', N'8 Chifley Square', N'AU', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'549300BUPYUQGB5BFX94', N'Hong Leong Bank', N'Jalan Dato Onn 50480 Kuala Lumpur', N'MY', 1)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'549300NROGNBV2T1GS07', N'Bank Negara Malaysia', N'Jalan Dato Onn', N'MY', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'984500653R409CC5AB28', N'JPM Chase', N'10 Shenton Wy', N'SG', 1)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'RVHOHKPBCJ2GSJ37YH94', N'Bank of Korea', N'39 Namdaemun-ro', N'KR', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'xxxxxxxxxxxxxxxxxxxx', N'Australia Commercial Bank', N'-------', N'AU', 1)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'yyyyyyyyyyyyyyyyyyyy', N'Korea Commercial Bank', N'-------', N'KR', 1)
GO
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'1', N'J.P. Admin', N'JPMadmin', N'admin', N'984500653R409CC5AB28')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'2', N'Monetary Authority of Singapore Admin', N'MASadmin', N'admin', N'54930035WQZLGC45RZ35')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'3', N'Hong Leong Bank Admin', N'HLBadmin', N'admin', N'549300BUPYUQGB5BFX94')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'4', N'Bank Negara Malaysia Admin', N'BNMadmin', N'admin', N'549300NROGNBV2T1GS07')
GO
SET IDENTITY_INSERT [dbo].[PolicyType] ON

INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (1, N'CFM', N'Capital Flow Management')
INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (2, N'SCL', N'Sanctions')
SET IDENTITY_INSERT [dbo].[PolicyType] OFF
GO
SET IDENTITY_INSERT [dbo].[Status] ON

INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (1, N'TransactionCreated')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (2, N'PoliciesApplied')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (3, N'ComplianceProofRequested')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (4, N'ComplianceCheckPassed')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (5, N'ProofInvalid')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (6, N'AssetSent')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (7, N'TransactionCompleted')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (8, N'TransactionCanceled')
SET IDENTITY_INSERT [dbo].[Status] OFF
GO
SET IDENTITY_INSERT [dbo].[TransactionType] ON

INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (1, N'DDWN', N'Loan Drawdown')
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (2, N'PPAY', N'Loan Repayment')
SET IDENTITY_INSERT [dbo].[TransactionType] OFF
GO
