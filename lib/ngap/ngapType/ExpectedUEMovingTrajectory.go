package ngapType

// Need to import "radio_simulator/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct ExpectedUEMovingTrajectory */
/* ExpectedUEMovingTrajectoryItem */
type ExpectedUEMovingTrajectory struct {
	List []ExpectedUEMovingTrajectoryItem `aper:"valueExt,sizeLB:1,sizeUB:16"`
}