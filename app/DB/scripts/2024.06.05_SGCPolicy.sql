USE [BIS1]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 2, 2, N'YYYYYYYYYYYYYYYYYYYY', N'SG', N'MY', N'SG', N'Public Sanctions List', 0, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
