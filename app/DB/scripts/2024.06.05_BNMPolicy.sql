USE [BIS4]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingCountryId], [OriginatingCountryId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 1, 1, 2, 1, N'100000000', 0, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
