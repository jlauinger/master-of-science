package eval2

type TotalPackageCounts struct {
	UnsafePointerTotal    	int
	UnsafeSizeofTotal		int
	UnsafeOffsetofTotal		int
	UnsafeAlignofTotal		int
	SliceHeaderTotal		int
	StringHeaderTotal		int
	UintptrTotal			int
}

type LocalPackageCounts struct {
	UnsafePointerLocal      int
	UnsafePointerVariable   int
	UnsafePointerParameter  int
	UnsafePointerAssignment int
	UnsafePointerCall       int
	UnsafePointerOther      int

	UnsafeSizeofLocal      int
	UnsafeSizeofVariable   int
	UnsafeSizeofParameter  int
	UnsafeSizeofAssignment int
	UnsafeSizeofCall       int
	UnsafeSizeofOther      int

	UnsafeOffsetofLocal      int
	UnsafeOffsetofVariable   int
	UnsafeOffsetofParameter  int
	UnsafeOffsetofAssignment int
	UnsafeOffsetofCall       int
	UnsafeOffsetofOther      int

	UnsafeAlignofLocal      int
	UnsafeAlignofVariable   int
	UnsafeAlignofParameter  int
	UnsafeAlignofAssignment int
	UnsafeAlignofCall       int
	UnsafeAlignofOther      int

	SliceHeaderLocal      int
	SliceHeaderVariable   int
	SliceHeaderParameter  int
	SliceHeaderAssignment int
	SliceHeaderCall       int
	SliceHeaderOther      int

	StringHeaderLocal      int
	StringHeaderVariable   int
	StringHeaderParameter  int
	StringHeaderAssignment int
	StringHeaderCall       int
	StringHeaderOther      int

	UintptrLocal      int
	UintptrVariable   int
	UintptrParameter  int
	UintptrAssignment int
	UintptrCall       int
	UintptrOther      int
}

