USE [BIS]
GO
/****** Object:  Table [dbo].[BankType]    Script Date: 24.04.2024. 23:20:31 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[BankType](
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](250) NOT NULL,
 CONSTRAINT [PK_BankType] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
SET IDENTITY_INSERT [dbo].[BankType] ON 

INSERT [dbo].[BankType] ([Id], [Name]) VALUES (1, N'Commercial')
INSERT [dbo].[BankType] ([Id], [Name]) VALUES (2, N'Central')
SET IDENTITY_INSERT [dbo].[BankType] OFF
GO
