USE [BIS2]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 2, 1, N'MY', N'SG', N'Public Sanctions List', 0, 1)
GO

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (2, 2, 1, N'MY', N'SG', N'Private Sanctions List', 1, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
