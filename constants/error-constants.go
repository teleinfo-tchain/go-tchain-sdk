/********************************************************************************
   This file is part of go-bif.
   go-bif is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-bif is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-bif.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package customerror

import "errors"

var (
	// EMPTYRESPONSE - Server response is empty
	EMPTYRESPONSE = errors.New("empty response")
	// UNPARSEABLEINTERFACE - the conversion failed
	UNPARSEABLEINTERFACE = errors.New("unParsable Interface")
	// WEBSOCKETNOTDENIFIED - Websocket connection dont exist
	WEBSOCKETNOTDENIFIED = errors.New("websocket connection dont exist")
)
