USE [BIS2]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 2, 1, N'XXXXXXXXXXXXXXXXXXXX', N'MY', N'SG', N'MY', N'Public Sanctions List', 0, 1)
GO

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [Owner], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [BeneficiaryJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (2, 2, 1, N'XXXXXXXXXXXXXXXXXXXX', N'MY', N'SG', N'MY', N'Private Sanctions List', 1, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
