USE [BIS2]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingCountryId], [OriginatingCountryId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 2, 1, 2, 1, N'Public Sanctions List', 0, 1)
GO

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingCountryId], [OriginatingCountryId], [Parameters], [IsPrivate], [Latest])
VALUES
    (2, 2, 1, 2, 1, N'Private Sanctions List', 1, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
