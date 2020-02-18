/*

Try catch functionality for areas where runtime error are very likely to occur.

*/

package CCMEException

type Catcher struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (TCF Catcher) Do() {
	if TCF.Finally != nil {
		defer TCF.Finally()
	}
	if TCF.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				TCF.Catch(r)
			}
		}()
	}
	TCF.Try()
}
