USE [BIS1]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 1, 1, N'yyyyyyyyyyyyyyyyyyyy', N'KR', N'KR', N'AU', N'Public - Open Sanctions', 0, 1)

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (2, 2, 1, N'yyyyyyyyyyyyyyyyyyyy', N'KR', N'KR', N'AU', N'5000,100000', 0, 1)

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (3, 3, 1, N'yyyyyyyyyyyyyyyyyyyy', N'KR', N'KR', N'AU', N'5000,10000', 0, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF