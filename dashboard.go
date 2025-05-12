package main

import (
	"strconv"
	"strings"
)

type instanceSelector struct {
	current *proxyInstance
}

func (is *instanceSelector) beforeTableHTML() string {
	instances, err := readRegistryEntries()
	sb := strings.Builder{}
	if err != nil {
		sb.WriteString("<mark>")
		sb.WriteString(err.Error())
		sb.WriteString("</mark>")
		return sb.String()
	}
	sb.WriteString("<select id=\"instance-selector\">")

	// Get current host:port
	currentHostPort := "localhost:" + strconv.Itoa(is.current.Port)

	for _, i := range instances {
		instanceURL := "http://" + i.Host + ":" + strconv.Itoa(i.Port)
		selected := ""
		if i.Host+":"+strconv.Itoa(i.Port) == currentHostPort {
			selected = " selected"
		}
		sb.WriteString("<option value=\"" + instanceURL + "\"" + selected + ">")
		if i.Host+":"+strconv.Itoa(i.Port) == currentHostPort {
			sb.WriteString("▼ " + i.Title + " :: " + i.Host + ":" + strconv.Itoa(i.Port) + " :: " + i.Command)
		} else {
			sb.WriteString(i.Title + " :: " + i.Host + ":" + strconv.Itoa(i.Port) + " :: " + i.Command)
		}
		sb.WriteString("</option>")
	}

	sb.WriteString("</select>")

	return sb.String()
}

func (is *instanceSelector) endHeadHTML() string {
	return `
			<script>
			document.addEventListener('DOMContentLoaded', function() {
				var selector = document.getElementById('instance-selector');
				if (selector) {
					selector.addEventListener('change', function() {
						window.location.href = this.value;
					});
				}
			});
			</script>
			`
}
