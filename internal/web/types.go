package web

import (
	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
	"github.com/dragon123098/Attendance-HackDay.git/internal/viewmodel"
)

type User = domain.User
type Classroom = domain.Classroom
type Schedule = domain.Schedule
type ShopItem = domain.ShopItem
type AvatarConfig = domain.AvatarConfig
type CoinTransaction = domain.CoinTransaction
type AttendanceRecord = domain.AttendanceRecord

type PageData = viewmodel.PageData
type ScheduleItemView = viewmodel.ScheduleItemView
type DoubleDayView = viewmodel.DoubleDayView
type ShopItemView = viewmodel.ShopItemView
type ThemeBackgroundOptionView = viewmodel.ThemeBackgroundOptionView
type AvatarBaseOptionView = viewmodel.AvatarBaseOptionView
type AvatarCosmeticOptionView = viewmodel.AvatarCosmeticOptionView
type AvatarLayerView = viewmodel.AvatarLayerView
type AvatarPreviewView = viewmodel.AvatarPreviewView
