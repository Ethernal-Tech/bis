USE [BIS4]
GO

SET IDENTITY_INSERT [dbo].[Policy] ON

INSERT [dbo].[Policy]
    ([Id], [PolicyTypeId], [TransactionTypeId], [PolicyEnforcingJurisdictionId], [OriginatingJurisdictionId], [Parameters], [IsPrivate], [Latest])
VALUES
    (1, 1, 1, N'MY', N'SG', N'100000000', 0, 1)
GO

SET IDENTITY_INSERT [dbo].[Policy] OFF
