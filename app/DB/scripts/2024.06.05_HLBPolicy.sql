USE [BIS2]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingCountryId], [OriginatingCountryId], [Parameters])
VALUES
    (1, 2, 1, 2, 1, N'Public Sanctions List')
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
